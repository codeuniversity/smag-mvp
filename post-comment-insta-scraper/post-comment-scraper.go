package post_comment_insta_scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/httpClient"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/segmentio/kafka-go"
	"time"
)

type PostCommentScraper struct {
	postIdQReader       *kafka.Reader
	commentsInfoQWriter *kafka.Writer
	errQWriter          *kafka.Writer
	*service.Executor
	kafkaAddress string
	httpClient   *httpClient.HttpClient
}

func New(kafkaAddress string) *PostCommentScraper {
	p := &PostCommentScraper{}
	p.postIdQReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "insta_comments_group1",
		Topic:          "user_post",
		CommitInterval: time.Minute * 40,
	})
	p.commentsInfoQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "insta_comments_info",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	p.errQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "post_comment_insta_scraper_errors",
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})
	p.Executor = service.New()
	p.kafkaAddress = kafkaAddress
	p.httpClient = httpClient.New(2, p.kafkaAddress)
	return p
}

func (p *PostCommentScraper) Run() {
	defer func() {
		p.MarkAsStopped()
	}()

	fmt.Println("starting Instagram post-scraper")
	counter := 0
	for p.IsRunning() {

		message, err := p.postIdQReader.FetchMessage(context.Background())

		if err != nil {
			fmt.Println(err)
			continue
		}

		var post models.InstagramPost
		err = json.Unmarshal(message.Value, &post)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var postsComments *models.InstaPostComments
		counter++
		err = p.httpClient.WithRetries(3, func() error {
			instaPostComments, err := p.httpClient.ScrapePostComments(post.ShortCode)

			if err != nil {
				return err
			}

			postsComments = &instaPostComments
			return nil
		})

		if err != nil {
			errorMessage := models.InstaCommentScrapError{
				PostId: post.PostId,
				Error:  err.Error(),
			}

			errorMessageJson, err := json.Marshal(errorMessage)
			if err != nil {
				panic(err)
			}
			p.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: errorMessageJson})
		} else {
			err = p.sendComments(postsComments, post)
			if err != nil {
				panic(err)
			}
		}
		p.postIdQReader.CommitMessages(context.Background(), message)
		counter++
	}
}

func (p *PostCommentScraper) sendComments(postsComments *models.InstaPostComments, postId models.InstagramPost) error {

	messages := make([]kafka.Message, 0, len(postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges))
	for _, element := range postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges {
		if element.Node.ID != "" {
			postComment := models.InstaComment{
				Id:        element.Node.ID,
				Text:      element.Node.Text,
				CreatedAt: element.Node.CreatedAt,
				PostId:    postId.PostId,
				ShortCode: postId.ShortCode,
				UserName:  element.Node.Owner.Username,
			}
			fmt.Println("CommentText: ", element.Node.Text)
			postCommentJson, err := json.Marshal(postComment)

			if err != nil {
				panic(fmt.Errorf("json marshal failed with InstaComment: %s", err))
			}

			m := kafka.Message{Value: postCommentJson}
			messages = append(messages, m)
		}
	}
	return p.commentsInfoQWriter.WriteMessages(context.Background(), messages...)
}

func (p *PostCommentScraper) Close() {
	p.Stop()
	p.WaitUntilStopped(time.Second * 3)

	p.postIdQReader.Close()
	p.commentsInfoQWriter.Close()
	p.errQWriter.Close()
	p.httpClient.Close()
	p.MarkAsClosed()
}
