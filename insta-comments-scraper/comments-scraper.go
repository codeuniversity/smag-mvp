package insta_comments_scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/scraper-client"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	userPostsCommentUrl = "https://www.instagram.com/graphql/query/?query_hash=865589822932d1b43dfe312121dd353a&variables=%s"
)

type PostCommentScraper struct {
	postIdQReader       *kafka.Reader
	commentsInfoQWriter *kafka.Writer
	errQWriter          *kafka.Writer
	*service.Executor
	httpClient scraper_client.ScraperClient
}

func New(postReader *kafka.Reader, commentsInfoQWriter *kafka.Writer, errQWriter *kafka.Writer) *PostCommentScraper {
	p := &PostCommentScraper{}
	p.postIdQReader = postReader
	p.commentsInfoQWriter = commentsInfoQWriter
	p.errQWriter = errQWriter
	p.Executor = service.New()
	p.httpClient = scraper_client.NewSimpleScraperClient()
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
		counter++
		postsComments, err := p.scrapeComments(post.ShortCode)

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

		time.Sleep(time.Millisecond * 900)
	}
}

func (p *PostCommentScraper) scrapeComments(shortCode string) (*instaPostComments, error) {
	var postsComments *instaPostComments
	err := p.httpClient.WithRetries(3, func() error {
		time.Sleep(1400 * time.Millisecond)
		instaPostComments, err := p.scrapePostComment(shortCode)

		if err != nil {
			return err
		}

		postsComments = &instaPostComments
		return nil
	})
	return postsComments, err
}

func (p *PostCommentScraper) sendComments(postsComments *instaPostComments, postId models.InstagramPost) error {
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

func (p *PostCommentScraper) scrapePostComment(shortCode string) (instaPostComments, error) {
	var instaPostComment instaPostComments
	type Variables struct {
		Shortcode           string `json:"shortcode"`
		ChildCommentCount   int    `json:"child_comment_count"`
		FetchCommentCount   int    `json:"fetch_comment_count"`
		ParentCommentCount  int    `json:"parent_comment_count"`
		HasThreadedComments bool   `json:"has_threaded_comments"`
	}

	variable := &Variables{shortCode, 3, 40, 24, true}
	variableJson, err := json.Marshal(variable)
	if err != nil {
		return instaPostComment, err
	}

	queryEncoded := url.QueryEscape(string(variableJson))
	url := fmt.Sprintf(userPostsCommentUrl, queryEncoded)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return instaPostComment, err
	}
	response, err := p.httpClient.Do(request)
	if err != nil {
		return instaPostComment, err
	}
	if response.StatusCode != 200 {
		return instaPostComment, &scraper_client.HttpStatusError{fmt.Sprintf("Error HttpStatus: %s", response.StatusCode)}
	}
	fmt.Println("ScrapePostComments got response")
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return instaPostComment, err
	}
	err = json.Unmarshal(body, &instaPostComment)
	if err != nil {
		return instaPostComment, err
	}
	return instaPostComment, nil
}

func (p *PostCommentScraper) Close() {
	p.Stop()
	p.WaitUntilStopped(time.Second * 3)

	p.postIdQReader.Close()
	p.commentsInfoQWriter.Close()
	p.errQWriter.Close()
	p.MarkAsClosed()
}
