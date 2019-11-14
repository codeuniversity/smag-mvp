package insta_likes_scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/models"
	client "github.com/codeuniversity/smag-mvp/scraper-client"
	"github.com/codeuniversity/smag-mvp/worker"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/segmentio/kafka-go"
)

const (
	//instaLikesInfoURL = "https://www.instagram.com/graphql/query/?query_hash=d5d763b1e2acf209d62d22d184488e57&variables=%s"
	instaLikesURL = "https://www.instagram.com/graphql/query/?query_hash=d5d763b1e2acf209d62d22d184488e57&variables=%s"
)

// PostLikesScraper scrapes the comments under post
type PostLikesScraper struct {
	*worker.Worker

	postIDQReader    *kafka.Reader
	likesInfoQWriter *kafka.Writer
	errQWriter       *kafka.Writer

	httpClient        client.ScraperClient
	requestRetryCount int
	LikesLimit        int
}

// New returns an initialized PostLikesScraper
func New(config *client.ScraperConfig, awsServiceAddress string, postIDQReader *kafka.Reader, commentsInfoQWriter *kafka.Writer, errQWriter *kafka.Writer, commentLimit int) *PostLikesScraper {
	s := &PostLikesScraper{}
	s.postIDQReader = postIDQReader
	s.likesInfoQWriter = commentsInfoQWriter
	s.errQWriter = errQWriter
	s.requestRetryCount = config.RequestRetryCount
	s.LikesLimit = commentLimit

	if awsServiceAddress == "" {
		s.httpClient = client.NewSimpleScraperClient()
	} else {
		s.httpClient = client.NewHttpClient(awsServiceAddress, config)
	}
	s.Worker = worker.Builder{}.WithName("insta_likes_scraper").
		WithWorkStep(s.runStep).
		AddShutdownHook("postIDQReader", postIDQReader.Close).
		AddShutdownHook("likesInfoQWriter", commentsInfoQWriter.Close).
		AddShutdownHook("errQWriter", errQWriter.Close).
		MustBuild()

	return s
}

func (s *PostLikesScraper) runStep() error {
	message, err := s.postIDQReader.FetchMessage(context.Background())

	log.Println("New Message")
	if err != nil {
		return err
	}

	var post models.InstagramPost
	err = json.Unmarshal(message.Value, &post)
	if err != nil {
		return err
	}

	log.Println("ShortCode: ", post.ShortCode)

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
		for (commentCounter < s.LikesLimit) && nextPage {
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
			commentCounter += 12
		}
	}
	if err != nil {
		return err
	}

	return s.postIDQReader.CommitMessages(context.Background(), message)
}

func (s *PostLikesScraper) sendInstaCommentError(postId string, err error) error {
	errorMessage := models.InstaCommentScrapError{
		PostID: postId,
		Error:  err.Error(),
	}

	errorMessageJSON, err := json.Marshal(errorMessage)
	if err != nil {
		panic(err)
	}
	return s.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: errorMessageJSON})
}

func (s *PostLikesScraper) scrapeCommentsInfo(shortCode string) (*instaPostCommentsInfo, error) {
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

func (s *PostLikesScraper) scrapeComments(shortCode string, cursor string) (*instaPostComments, error) {
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

func (s *PostLikesScraper) sendKafkaComments(postsComments *instaPostCommentsInfo, postID models.InstagramPost) error {
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
	return s.likesInfoQWriter.WriteMessages(context.Background(), messages...)
}

func (s *PostLikesScraper) sendCommentsInfoRequest(shortCode string) (InstaPostLikes, error) {
	var instaPostComment InstaPostLikes
	type Variables struct {
		Shortcode   string `json:"shortcode"`
		IncludeReel bool   `json:"include_reel"`
		First       int    `json:"first"`
	}

	variable := &Variables{shortCode, true, 24}
	variableJSON, err := json.Marshal(variable)
	if err != nil {
		return instaPostComment, err
	}

	queryEncoded := url.QueryEscape(string(variableJSON))
	url := fmt.Sprintf(instaLikesURL, queryEncoded)

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

func (s *PostLikesScraper) sendCommentsRequest(shortCode string, cursor string) (InstaPostLikes, error) {

	var instaPostComment InstaPostLikes
	type Variables struct {
		Shortcode   string `json:"shortcode"`
		IncludeReel bool   `json:"include_reel"`
		First       int    `json:"first"`
		After       string `json:"after"`
	}

	variable := &Variables{shortCode, true, 12, cursor}

	variableJSON, err := json.Marshal(variable)
	if err != nil {
		return instaPostComment, err
	}

	queryEncoded := url.QueryEscape(string(variableJSON))
	url := fmt.Sprintf(instaLikesURL, queryEncoded)

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
