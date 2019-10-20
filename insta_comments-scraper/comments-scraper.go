package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/codeuniversity/smag-mvp/models"
	client "github.com/codeuniversity/smag-mvp/scraper-client"

	"github.com/codeuniversity/smag-mvp/service"
	"github.com/segmentio/kafka-go"
)

const (
	userPostsCommentURL = "https://www.instagram.com/graphql/query/?query_hash=865589822932d1b43dfe312121dd353a&variables=%s"
)

// PostCommentScraper scrapes the comments under post
type PostCommentScraper struct {
	postIDQReader       *kafka.Reader
	commentsInfoQWriter *kafka.Writer
	errQWriter          *kafka.Writer
	*service.Executor
	httpClient client.ScraperClient
}

// New returns an initialized PostCommentScraper
func New(postReader *kafka.Reader, commentsInfoQWriter *kafka.Writer, errQWriter *kafka.Writer) *PostCommentScraper {
	p := &PostCommentScraper{}
	p.postIDQReader = postReader
	p.commentsInfoQWriter = commentsInfoQWriter
	p.errQWriter = errQWriter
	p.Executor = service.New()
	p.httpClient = client.NewSimpleScraperClient()
	return p
}

// Run ...
func (p *PostCommentScraper) Run() {
	defer func() {
		p.MarkAsStopped()
	}()

	fmt.Println("starting Instagram post-scraper")
	counter := 0
	for p.IsRunning() {

		message, err := p.postIDQReader.FetchMessage(context.Background())

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
				PostID: post.PostID,
				Error:  err.Error(),
			}

			errorMessageJSON, err := json.Marshal(errorMessage)
			if err != nil {
				panic(err)
			}
			p.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: errorMessageJSON})
		} else {
			err = p.sendComments(postsComments, post)
			if err != nil {
				panic(err)
			}
		}
		p.postIDQReader.CommitMessages(context.Background(), message)
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

func (p *PostCommentScraper) sendComments(postsComments *instaPostComments, postID models.InstagramPost) error {
	fmt.Println("sendComments: ", len(postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges))
	messages := make([]kafka.Message, 0, len(postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges))
	for _, element := range postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges {
		if element.Node.ID != "" {
			postComment := models.InstaComment{
				ID:            element.Node.ID,
				Text:          element.Node.Text,
				CreatedAt:     element.Node.CreatedAt,
				PostID:        postID.PostID,
				ShortCode:     postID.ShortCode,
				OwnerUsername: element.Node.Owner.Username,
			}
			fmt.Println("CommentText: ", element.Node.Text)
			postCommentJSON, err := json.Marshal(postComment)

			if err != nil {
				panic(fmt.Errorf("json marshal failed with InstaComment: %s", err))
			}

			m := kafka.Message{Value: postCommentJSON}
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
	variableJSON, err := json.Marshal(variable)
	if err != nil {
		return instaPostComment, err
	}

	queryEncoded := url.QueryEscape(string(variableJSON))
	url := fmt.Sprintf(userPostsCommentURL, queryEncoded)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return instaPostComment, err
	}
	response, err := p.httpClient.Do(request)
	if err != nil {
		return instaPostComment, err
	}
	if response.StatusCode != 200 {
		return instaPostComment, &client.HTTPStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
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

// Close ...
func (p *PostCommentScraper) Close() {
	p.Stop()
	p.WaitUntilStopped(time.Second * 3)

	p.postIDQReader.Close()
	p.commentsInfoQWriter.Close()
	p.errQWriter.Close()
	p.MarkAsClosed()
}
