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
	scraper_client "github.com/codeuniversity/smag-mvp/scraper-client"
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

	requestCounter int

	httpClient scraper_client.ScraperClient
}

// New returns an initilized scraper
func New(nameQReader *kafka.Reader, infoQWriter *kafka.Writer, errQWriter *kafka.Writer) *InstaPostsScraper {
	i := &InstaPostsScraper{}
	i.nameQReader = nameQReader
	i.postsQWriter = infoQWriter
	i.errQWriter = errQWriter

	i.httpClient = scraper_client.NewSimpleScraperClient()

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
	fmt.Println("fetching")
	m, err := i.nameQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	username := string(m.Value)

	instagramAccountInfo, err := i.accountInfo(username)
	i.requestCounter++
	fmt.Println("Instagram Request Counter: ", i.requestCounter)
	if err != nil {
		errMessage := &models.InstagramScrapeError{
			Name:  username,
			Error: err.Error(),
		}
		serializedErr, err := json.Marshal(errMessage)
		if err != nil {
			return err
		}
		i.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
		i.nameQReader.CommitMessages(context.Background(), m)
		return nil
	}

	if instagramAccountInfo.Graphql.User.IsPrivate {
		fmt.Println("Username: ", username, " is private")
		i.nameQReader.CommitMessages(context.Background(), m)
		return nil
	}
	fmt.Println("Username: ", username, " Posts Init")
	cursor := instagramAccountInfo.Graphql.User.EdgeOwnerToTimelineMedia.PageInfo.EndCursor
	userID := instagramAccountInfo.Graphql.User.ID

	err = i.sendUserInfoPostsID(instagramAccountInfo, username, userID)
	if err != nil {
		return err
	}

	isPostsSendingFinished := false
	for isPostsSendingFinished {
		fmt.Println("Username: ", username, " accountPosts")
		accountMedia, err := i.accountPosts(userID, cursor)
		i.requestCounter++
		fmt.Println("Instagram Request Counter: ", i.requestCounter)

		if err != nil {
			return i.sendErrorMessage(m, username, err)
		}

		cursor = accountMedia.Data.User.EdgeOwnerToTimelineMedia.PageInfo.EndCursor
		i.sendUserTimlinePostsID(accountMedia, username, userID)

		if !accountMedia.Data.User.EdgeOwnerToTimelineMedia.PageInfo.HasNextPage {
			isPostsSendingFinished = false
		}
		time.Sleep(time.Millisecond * 100)
	}
	return i.nameQReader.CommitMessages(context.Background(), m)
}

func (i *InstaPostsScraper) accountInfo(username string) (*instagramAccountInfo, error) {
	var instagramAccountInfo *instagramAccountInfo

	err := i.httpClient.WithRetries(2, func() error {
		accountInfo, err := i.scrapeAccountInfo(username)
		if err != nil {
			return err
		}

		instagramAccountInfo = &accountInfo
		return nil
	})

	if err != nil {
		return instagramAccountInfo, err
	}
	return instagramAccountInfo, err
}

func (i *InstaPostsScraper) accountPosts(userID string, cursor string) (*instagramMedia, error) {
	var instagramAccountMedia *instagramMedia

	err := i.httpClient.WithRetries(2, func() error {
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
		return userAccountInfo, &scraper_client.HTTPStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
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
	fmt.Println(string(variableJSON))
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
		return media, &scraper_client.HTTPStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
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

func (i *InstaPostsScraper) sendUserInfoPostsID(instagramAccountInfo *instagramAccountInfo, username string, userID string) error {
	for _, element := range instagramAccountInfo.Graphql.User.EdgeOwnerToTimelineMedia.Edges {
		fmt.Println("Edges ", username)
		fmt.Println(element.Node.Typename)
		if element.Node.Shortcode != "" {
			fmt.Println("Edges Node1 ", username)
			instagramPost := models.InstagramPost{
				PostID:     element.Node.ID,
				ShortCode:  element.Node.Shortcode,
				UserID:     userID,
				UserName:   username,
				PictureURL: element.Node.DisplayURL}

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

func (i *InstaPostsScraper) sendUserTimlinePostsID(accountMedia *instagramMedia, username string, userID string) {
	for _, element := range accountMedia.Data.User.EdgeOwnerToTimelineMedia.Edges {
		if element.Node.ID != "" {

			var postsTaggedUsers []string
			for _, element := range element.Node.EdgeMediaToTaggedUser.Edges {
				postsTaggedUsers = append(postsTaggedUsers, element.Node.User.Username)
			}

			instagramPost := models.InstagramPost{
				PostID:      element.Node.ID,
				ShortCode:   element.Node.Shortcode,
				UserID:      userID,
				UserName:    username,
				PictureURL:  element.Node.DisplayURL,
				TaggedUsers: postsTaggedUsers,
			}

			instagramPostJSON, err := json.Marshal(instagramPost)

			if err != nil {
				fmt.Println(err)
				break
			}

			err = i.postsQWriter.WriteMessages(context.Background(), kafka.Message{Value: instagramPostJSON})
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
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
