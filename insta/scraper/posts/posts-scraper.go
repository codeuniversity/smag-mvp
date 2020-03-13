package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/codeuniversity/smag-mvp/insta/models"
	client "github.com/codeuniversity/smag-mvp/scraper-client"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/segmentio/kafka-go"
)

const (
	userAccountInfoURL  = "https://instagram.com/%s/?__a=1"
	userAccountMediaURL = "https://www.instagram.com/graphql/query/?query_hash=58b6785bea111c67129decbe6a448951&variables=%s"
)

// InstaPostsScraper scrapes posts from instagram
type InstaPostsScraper struct {
	*worker.Worker

	nameQReader  *kafka.Reader
	postsQWriter *kafka.Writer
	errQWriter   *kafka.Writer

	requestCounter    int
	httpClient        client.ScraperClient
	requestRetryCount int
}

// New returns an initilized scraper
func New(config *client.ScraperConfig, awsServiceAddress string, nameQReader *kafka.Reader, infoQWriter *kafka.Writer, errQWriter *kafka.Writer) *InstaPostsScraper {
	i := &InstaPostsScraper{}
	i.nameQReader = nameQReader
	i.postsQWriter = infoQWriter
	i.errQWriter = errQWriter
	i.requestRetryCount = config.RequestRetryCount

	if awsServiceAddress == "" {
		i.httpClient = client.NewSimpleScraperClient()
	} else {
		i.httpClient = client.NewHttpClient(awsServiceAddress, config)
	}

	i.Worker = worker.Builder{}.WithName("insta_posts_scraper").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("nameQReader", nameQReader.Close).
		AddShutdownHook("infoQWriter", infoQWriter.Close).
		AddShutdownHook("errQWriter", errQWriter.Close).
		MustBuild()

	return i
}

func (i *InstaPostsScraper) runStep() error {
	log.Println("fetching")
	m, err := i.nameQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	username := string(m.Value)

	instagramAccountInfo, err := i.accountInfo(username)
	i.requestCounter++
	log.Println("Instagram Request Counter: ", i.requestCounter)
	if err != nil {
		errMessage := &models.InstagramScrapeError{
			Name:  username,
			Error: err.Error(),
		}
		serializedErr, err := json.Marshal(errMessage)
		if err != nil {
			return err
		}
		err = i.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
		if err != nil {
			return err
		}
		return i.nameQReader.CommitMessages(context.Background(), m)
	}

	if instagramAccountInfo == nil {
		log.Println("InstagramAccount is nil")

		errMessage := &models.InstagramScrapeError{
			Name: username,
		}
		serializedErr, err := json.Marshal(errMessage)
		if err != nil {
			return err
		}
		err = i.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})

		if err != nil {
			return err
		}
		return i.nameQReader.CommitMessages(context.Background(), m)
	}

	if instagramAccountInfo.Graphql.User.IsPrivate {
		log.Println("Username: ", username, " is private")
		return i.nameQReader.CommitMessages(context.Background(), m)
	}
	log.Println("Username: ", username, " Posts Init")
	userID := instagramAccountInfo.Graphql.User.ID

	isPostsSendingFinished := false
	cursor := ""
	for !isPostsSendingFinished {
		log.Println("Username: ", username, " accountPosts")
		accountMedia, err := i.accountPosts(userID, cursor)
		i.requestCounter++
		log.Println("Instagram Request Counter: ", i.requestCounter)

		if err != nil {
			return i.sendErrorMessage(m, username, err)
		}

		cursor = accountMedia.Data.User.EdgeOwnerToTimelineMedia.PageInfo.EndCursor
		err = i.sendUserTimlinePostsID(accountMedia, username, userID)

		if err != nil {
			return err
		}
		if !accountMedia.Data.User.EdgeOwnerToTimelineMedia.PageInfo.HasNextPage {
			isPostsSendingFinished = true
		}
	}
	return i.nameQReader.CommitMessages(context.Background(), m)
}

func (i *InstaPostsScraper) accountInfo(username string) (*instagramAccountInfo, error) {
	var instagramAccountInfo *instagramAccountInfo

	err := i.httpClient.WithRetries(i.requestRetryCount, func() error {
		accountInfo, err := i.scrapeAccountInfo(username)
		if err != nil {
			return err
		}

		instagramAccountInfo = &accountInfo
		return nil
	})

	if err != nil {
		return nil, err
	}
	return instagramAccountInfo, err
}

func (i *InstaPostsScraper) accountPosts(userID string, cursor string) (*instagramMedia, error) {
	var instagramAccountMedia *instagramMedia

	err := i.httpClient.WithRetries(i.requestRetryCount, func() error {
		accountInfo, err := i.scrapeProfileMedia(userID, cursor)
		if err != nil {
			return err
		}

		instagramAccountMedia = &accountInfo
		return nil
	})

	if err != nil {
		return instagramAccountMedia, err
	}
	return instagramAccountMedia, err
}

func (i *InstaPostsScraper) scrapeAccountInfo(username string) (instagramAccountInfo, error) {
	var userAccountInfo instagramAccountInfo
	url := fmt.Sprintf(userAccountInfoURL, username)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return userAccountInfo, err
	}

	response, err := i.httpClient.Do(request)
	if err != nil {
		return userAccountInfo, err
	}
	if response.StatusCode != 200 {
		return userAccountInfo, &client.HTTPStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return userAccountInfo, err
	}
	err = json.Unmarshal(body, &userAccountInfo)
	if err != nil {
		return userAccountInfo, err
	}
	return userAccountInfo, nil
}

func (i *InstaPostsScraper) scrapeProfileMedia(userID string, endCursor string) (instagramMedia, error) {
	var media instagramMedia

	type Variables struct {
		ID    string `json:"id"`
		First int    `json:"first"`
		After string `json:"after"`
	}
	variable := &Variables{userID, 12, endCursor}
	variableJSON, err := json.Marshal(variable)
	if err != nil {
		return media, err
	}
	queryEncoded := url.QueryEscape(string(variableJSON))
	url := fmt.Sprintf(userAccountMediaURL, queryEncoded)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return media, err
	}
	response, err := i.httpClient.Do(request)
	if err != nil {
		return media, err
	}
	if response.StatusCode != 200 {
		return media, &client.HTTPStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return media, err
	}
	err = json.Unmarshal(body, &media)
	if err != nil {
		return media, err
	}
	return media, nil
}

func (i *InstaPostsScraper) sendUserTimlinePostsID(accountMedia *instagramMedia, username string, userID string) error {
	for _, element := range accountMedia.Data.User.EdgeOwnerToTimelineMedia.Edges {
		if element.Node.ID != "" {

			var postsTaggedUsers []string
			for _, element := range element.Node.EdgeMediaToTaggedUser.Edges {
				postsTaggedUsers = append(postsTaggedUsers, element.Node.User.Username)
			}

			var caption string
			if len(element.Node.EdgeMediaToCaption.Edges) > 0 {
				caption = element.Node.EdgeMediaToCaption.Edges[0].Node.Text
			}

			instagramPost := models.InstagramPost{
				PostID:      element.Node.ID,
				ShortCode:   element.Node.Shortcode,
				UserID:      userID,
				UserName:    username,
				PictureURL:  element.Node.DisplayURL,
				TaggedUsers: postsTaggedUsers,
				Caption:     caption,
			}

			instagramPostJSON, err := json.Marshal(instagramPost)

			if err != nil {
				return err
			}

			err = i.postsQWriter.WriteMessages(context.Background(), kafka.Message{Value: instagramPostJSON})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *InstaPostsScraper) sendErrorMessage(m kafka.Message, username string, errToSend error) error {
	errMessage := &models.InstagramScrapeError{
		Name:  username,
		Error: errToSend.Error(),
	}
	serializedErr, err := json.Marshal(errMessage)
	if err != nil {
		return err
	}
	err = i.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
	if err != nil {
		return err
	}
	err = i.nameQReader.CommitMessages(context.Background(), m)
	if err != nil {
		return err
	}
	return nil
}
