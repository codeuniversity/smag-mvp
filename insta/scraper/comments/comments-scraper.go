package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/codeuniversity/smag-mvp/insta/models"
	client "github.com/codeuniversity/smag-mvp/scraper-client"
	"github.com/codeuniversity/smag-mvp/worker"

	"github.com/segmentio/kafka-go"
)

const (
	instaCommentInfoURL = "https://www.instagram.com/graphql/query/?query_hash=fead941d698dc1160a298ba7bec277ac&variables=%s"
	instaCommentURL     = "https://www.instagram.com/graphql/query/?query_hash=97b41c52301f77ce508f55e66d17620e&variables=%s"
)

const instaCommentCount = 12

// PostCommentScraper scrapes the comments under post
type PostCommentScraper struct {
	*worker.Worker

	postIDQReader       *kafka.Reader
	commentsInfoQWriter *kafka.Writer
	errQWriter          *kafka.Writer

	httpClient        client.ScraperClient
	requestRetryCount int
	commentLimit      int
}

// New returns an initialized PostCommentScraper
func New(config *client.ScraperConfig, awsServiceAddress string, postIDQReader *kafka.Reader, commentsInfoQWriter *kafka.Writer, errQWriter *kafka.Writer, commentLimit int) *PostCommentScraper {
	s := &PostCommentScraper{}
	s.postIDQReader = postIDQReader
	s.commentsInfoQWriter = commentsInfoQWriter
	s.errQWriter = errQWriter
	s.requestRetryCount = config.RequestRetryCount
	s.commentLimit = commentLimit

	if awsServiceAddress == "" {
		s.httpClient = client.NewSimpleScraperClient()
	} else {
		s.httpClient = client.NewHttpClient(awsServiceAddress, config)
	}
	s.Worker = worker.Builder{}.WithName("insta_comments_scraper").
		WithWorkStep(s.runStep).
		AddShutdownHook("postIDQReader", postIDQReader.Close).
		AddShutdownHook("commentsInfoQWriter", commentsInfoQWriter.Close).
		AddShutdownHook("errQWriter", errQWriter.Close).
		MustBuild()

	return s
}

func (s *PostCommentScraper) runStep() error {
	message, err := s.postIDQReader.FetchMessage(context.Background())

	if err != nil {
		return err
	}

	var post models.InstagramPost
	err = json.Unmarshal(message.Value, &post)
	if err != nil {
		return err
	}

	postsComments, err := s.scrapeCommentsInfo(post.ShortCode)
	if err != nil {
		err := s.sendInstaCommentError(post.PostID, err)
		if err != nil {
			return err
		}
	} else {
		err = s.sendKafkaComments(postsComments, post)
		if err != nil {
			return err
		}

		endcursor := postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.PageInfo.EndCursor
		nextPage := postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.PageInfo.HasNextPage

		commentCounter := 24
		for (commentCounter < s.commentLimit) && nextPage {
			postComments, err := s.scrapeComments(post.ShortCode, endcursor)

			if err != nil {
				err := s.sendInstaCommentError(post.PostID, err)
				if err != nil {
					return err
				}
				continue
			}
			err = s.sendKafkaComments(postsComments, post)
			if err != nil {
				return err
			}

			nextPage = postComments.Data.ShortcodeMedia.EdgeMediaToParentComment.PageInfo.HasNextPage
			endcursor = postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.PageInfo.EndCursor
			commentCounter += instaCommentCount
		}
	}
	if err != nil {
		return err
	}

	return s.postIDQReader.CommitMessages(context.Background(), message)
}

func (s *PostCommentScraper) sendInstaCommentError(postId string, err error) error {
	errorMessage := models.InstaPostScrapeError{
		PostID: postId,
		Error:  err.Error(),
	}

	errorMessageJSON, err := json.Marshal(errorMessage)
	if err != nil {
		panic(err)
	}
	return s.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: errorMessageJSON})
}

func (s *PostCommentScraper) scrapeCommentsInfo(shortCode string) (*instaPostCommentsInfo, error) {
	var postsComments *instaPostCommentsInfo
	err := s.httpClient.WithRetries(s.requestRetryCount, func() error {
		instaPostComments, err := s.sendCommentsInfoRequest(shortCode)

		if err != nil {
			return err
		}

		postsComments = &instaPostComments
		return nil
	})
	return postsComments, err
}

func (s *PostCommentScraper) scrapeComments(shortCode string, cursor string) (*instaPostComments, error) {
	var postsComments *instaPostComments
	err := s.httpClient.WithRetries(s.requestRetryCount, func() error {
		instaPostComments, err := s.sendCommentsRequest(shortCode, cursor)

		if err != nil {
			return err
		}

		postsComments = &instaPostComments
		return nil
	})
	return postsComments, err
}

func (s *PostCommentScraper) sendKafkaComments(postsComments *instaPostCommentsInfo, postID models.InstagramPost) error {
	if len(postsComments.Data.ShortcodeMedia.EdgeMediaToParentComment.Edges) == 0 {
		return nil
	}
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
			log.Println("CommentText: ", element.Node.Text)
			postCommentJSON, err := json.Marshal(postComment)

			if err != nil {
				panic(fmt.Errorf("json marshal failed with InstaComment: %s", err))
			}

			m := kafka.Message{Value: postCommentJSON}
			messages = append(messages, m)
		}
	}
	return s.commentsInfoQWriter.WriteMessages(context.Background(), messages...)
}

func (s *PostCommentScraper) sendCommentsInfoRequest(shortCode string) (instaPostCommentsInfo, error) {
	var instaPostComment instaPostCommentsInfo
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
	url := fmt.Sprintf(instaCommentInfoURL, queryEncoded)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return instaPostComment, err
	}
	response, err := s.httpClient.Do(request)
	if err != nil {
		return instaPostComment, err
	}
	if response.StatusCode != 200 {
		return instaPostComment, &client.HTTPStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
	}
	log.Println("ScrapePostComments got response")
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

func (s *PostCommentScraper) sendCommentsRequest(shortCode string, cursor string) (instaPostComments, error) {

	var instaPostComment instaPostComments
	type Variables struct {
		Shortcode string `json:"shortcode"`
		First     int    `json:"first"`
		After     string `json:"after"`
	}

	variable := &Variables{shortCode, instaCommentCount, cursor}

	variableJSON, err := json.Marshal(variable)
	if err != nil {
		return instaPostComment, err
	}

	queryEncoded := url.QueryEscape(string(variableJSON))
	url := fmt.Sprintf(instaCommentURL, queryEncoded)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return instaPostComment, err
	}
	response, err := s.httpClient.Do(request)
	if err != nil {
		return instaPostComment, err
	}
	if response.StatusCode != 200 {
		return instaPostComment, &client.HTTPStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
	}
	log.Println("ScrapePostComments got response")
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
