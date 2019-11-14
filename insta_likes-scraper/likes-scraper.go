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

// PostLikesScraper scrapes the likes
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
func New(config *client.ScraperConfig, awsServiceAddress string, postIDQReader *kafka.Reader, likesInfoQWriter *kafka.Writer, errQWriter *kafka.Writer, likesLimit int) *PostLikesScraper {
	s := &PostLikesScraper{}
	s.postIDQReader = postIDQReader
	s.likesInfoQWriter = likesInfoQWriter
	s.errQWriter = errQWriter
	s.requestRetryCount = config.RequestRetryCount
	s.LikesLimit = likesLimit

	if awsServiceAddress == "" {
		s.httpClient = client.NewSimpleScraperClient()
	} else {
		s.httpClient = client.NewHttpClient(awsServiceAddress, config)
	}
	s.Worker = worker.Builder{}.WithName("insta_likes_scraper").
		WithWorkStep(s.runStep).
		AddShutdownHook("postIDQReader", postIDQReader.Close).
		AddShutdownHook("likesInfoQWriter", likesInfoQWriter.Close).
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

	postsComments, err := s.scrapeLikesInfo(post.ShortCode)
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

		endcursor := postsComments.Data.ShortcodeMedia.EdgeLikedBy.PageInfo.EndCursor
		nextPage := postsComments.Data.ShortcodeMedia.EdgeLikedBy.PageInfo.HasNextPage

		likeCounter := 24
		for (likeCounter < s.LikesLimit) && nextPage {
			postComments, err := s.scrapeLikes(post.ShortCode, endcursor)

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

			nextPage = postComments.Data.ShortcodeMedia.EdgeLikedBy.PageInfo.HasNextPage
			endcursor = postsComments.Data.ShortcodeMedia.EdgeLikedBy.PageInfo.EndCursor
			likeCounter += 12
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

func (s *PostLikesScraper) scrapeLikesInfo(shortCode string) (*InstaPostLikes, error) {
	var postsComments *InstaPostLikes
	err := s.httpClient.WithRetries(s.requestRetryCount, func() error {
		instaPostComments, err := s.sendLikesInfoRequest(shortCode)

		if err != nil {
			return err
		}

		postsComments = &instaPostComments
		return nil
	})
	return postsComments, err
}

func (s *PostLikesScraper) scrapeLikes(shortCode string, cursor string) (*InstaPostLikes, error) {
	var postsComments *InstaPostLikes
	err := s.httpClient.WithRetries(s.requestRetryCount, func() error {
		instaPostComments, err := s.sendLikesRequest(shortCode, cursor)

		if err != nil {
			return err
		}

		postsComments = &instaPostComments
		return nil
	})
	return postsComments, err
}

func (s *PostLikesScraper) sendKafkaComments(postsLikes *InstaPostLikes, postID models.InstagramPost) error {
	if len(postsLikes.Data.ShortcodeMedia.EdgeLikedBy.Edges) == 0 {
		return nil
	}
	messages := make([]kafka.Message, 0, len(postsLikes.Data.ShortcodeMedia.EdgeLikedBy.Edges))
	for _, element := range postsLikes.Data.ShortcodeMedia.EdgeLikedBy.Edges {
		if element.Node.ID != "" {
			postComment := models.InstaLikes{
				ID:            element.Node.ID,
				PostID:        postID.PostID,
				OwnerUsername: element.Node.Username,
				RealName:      element.Node.FullName,
				AvatarURL:     element.Node.ProfilePicURL,
			}
			log.Println("Likes: ", element.Node.Username)
			postCommentJSON, err := json.Marshal(postComment)

			if err != nil {
				panic(fmt.Errorf("json marshal failed with InstaLikes: %s", err))
			}

			m := kafka.Message{Value: postCommentJSON}
			messages = append(messages, m)
		}
	}
	return s.likesInfoQWriter.WriteMessages(context.Background(), messages...)
}

func (s *PostLikesScraper) sendLikesInfoRequest(shortCode string) (InstaPostLikes, error) {
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

func (s *PostLikesScraper) sendLikesRequest(shortCode string, cursor string) (InstaPostLikes, error) {

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
