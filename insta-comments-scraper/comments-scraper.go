package insta_comments_scraper

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
	httpClient *httpClient.HttpClient
}

func New(postReader *kafka.Reader, commentsInfoQWriter *kafka.Writer, errQWriter *kafka.Writer) *PostCommentScraper {
	p := &PostCommentScraper{}
	p.postIdQReader = postReader
	p.commentsInfoQWriter = commentsInfoQWriter
	p.errQWriter = errQWriter
	p.Executor = service.New()
	p.httpClient = httpClient.New(2, "52.58.171.160:9092")
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

		fmt.Println("New Message")
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

		fmt.Println("ShortCode: ", post.ShortCode)
		var postsComments *models.InstaPostComments
		counter++
		err = p.httpClient.WithRetries(3, func() error {
			time.Sleep(1400 * time.Millisecond)
			instaPostComments, err := p.httpClient.ScrapePostComments(post.ShortCode)

			if err != nil {
				return err
			}

			postsComments = &instaPostComments
			return nil
		})

		if err != nil {
			fmt.Println("error: ", err)
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

		time.Sleep(time.Millisecond * 900)
	}
}

func (p *PostCommentScraper) Scrape(shortCode string) (*models.InstaPostComments, error) {
	var postsComments *models.InstaPostComments
	err := p.httpClient.WithRetries(3, func() error {
		time.Sleep(9 * time.Second)
		instaPostComments, err := p.httpClient.ScrapePostComments(shortCode)

		if err != nil {
			return err
		}

		postsComments = &instaPostComments
		return nil
	})
	return postsComments, err
}

func (p *PostCommentScraper) sendComments(postsComments *models.InstaPostComments, postId models.InstagramPost) error {
	fmt.Println("sendComments: ", len(postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges))
	messages := make([]kafka.Message, 0, len(postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges))
	for _, element := range postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges {
		if element.Node.ID != "" {
			postComment := models.InstaComment{
				Id:            element.Node.ID,
				Text:          element.Node.Text,
				CreatedAt:     element.Node.CreatedAt,
				PostId:        postId.PostId,
				ShortCode:     postId.ShortCode,
				OwnerUsername: element.Node.Owner.Username,
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
