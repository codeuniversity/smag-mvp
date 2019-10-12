package insta_posts_scraper

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
	userAccountInfoUrl  = "https://instagram.com/%s/?__a=1"
	userAccountMediaUrl = "https://www.instagram.com/graphql/query/?query_hash=58b6785bea111c67129decbe6a448951&variables=%s"
)

type InstaPostsScraper struct {
	nameQReader  *kafka.Reader
	postsQWriter *kafka.Writer
	errQWriter   *kafka.Writer
	*service.Executor
	kafkaAddress string
	httpClient   scraper_client.ScraperClient
}

// New returns an initilized scraper
func New(nameQReader *kafka.Reader, infoQWriter *kafka.Writer, errQWriter *kafka.Writer) *InstaPostsScraper {
	i := &InstaPostsScraper{}
	i.nameQReader = nameQReader
	i.postsQWriter = infoQWriter
	i.errQWriter = errQWriter
	i.Executor = service.New()
	i.httpClient = scraper_client.NewSimpleScraperClient()
	return i
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

func (i *InstaPostsScraper) accountPosts(userId string, cursor string) (*InstagramMedia, error) {
	var instagramAccountMedia *InstagramMedia

	err := i.httpClient.WithRetries(2, func() error {
		accountInfo, err := i.scrapeProfileMedia(userId, cursor)
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
	url := fmt.Sprintf(userAccountInfoUrl, username)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return userAccountInfo, err
	}

	response, err := i.httpClient.Do(request)
	if err != nil {
		return userAccountInfo, err
	}
	if response.StatusCode != 200 {
		return userAccountInfo, &scraper_client.HttpStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
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

func (i *InstaPostsScraper) scrapeProfileMedia(userId string, endCursor string) (InstagramMedia, error) {
	var instagramMedia InstagramMedia

	type Variables struct {
		Id    string `json:"id"`
		First int    `json:"first"`
		After string `json:"after"`
	}
	variable := &Variables{userId, 12, endCursor}
	variableJson, err := json.Marshal(variable)
	fmt.Println(string(variableJson))
	if err != nil {
		return instagramMedia, err
	}
	queryEncoded := url.QueryEscape(string(variableJson))
	url := fmt.Sprintf(userAccountMediaUrl, queryEncoded)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return instagramMedia, err
	}
	response, err := i.httpClient.Do(request)
	if err != nil {
		return instagramMedia, err
	}
	if response.StatusCode != 200 {
		return instagramMedia, &scraper_client.HttpStatusError{S: fmt.Sprintf("Error HttpStatus: %d", response.StatusCode)}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return instagramMedia, err
	}
	err = json.Unmarshal(body, &instagramMedia)
	if err != nil {
		return instagramMedia, err
	}
	return instagramMedia, nil
}

func (i *InstaPostsScraper) Run() {
	defer func() {
		i.MarkAsStopped()
	}()

	fmt.Println("starting Instagram post-scraper")
	counter := 0
	for i.IsRunning() {
		fmt.Println("fetching")
		m, err := i.nameQReader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}

		username := string(m.Value)

		instagramAccountInfo, err := i.accountInfo(username)
		counter++
		fmt.Println("Instagram Request Counter: ", counter)
		if err != nil {
			errMessage := &models.InstagramScrapeError{
				Name:  username,
				Error: err.Error(),
			}
			serializedErr, err := json.Marshal(errMessage)
			if err != nil {
				fmt.Println(err)
				break
			}
			i.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
			i.nameQReader.CommitMessages(context.Background(), m)
			continue
		}

		if instagramAccountInfo.Graphql.User.IsPrivate {
			fmt.Println("Username: ", username, " is private")
			i.nameQReader.CommitMessages(context.Background(), m)
			continue
		}
		fmt.Println("Username: ", username, " Posts Init")
		cursor := instagramAccountInfo.Graphql.User.EdgeOwnerToTimelineMedia.PageInfo.EndCursor
		userId := instagramAccountInfo.Graphql.User.ID

		go i.sendUserInfoPostsId(instagramAccountInfo, username, userId)

		isPostsSendingFinished := false
		for i.IsRunning() && isPostsSendingFinished {
			fmt.Println("Username: ", username, " accountPosts")
			accountMedia, err := i.accountPosts(userId, cursor)
			counter++
			fmt.Println("Instagram Request Counter: ", counter)

			if err != nil {
				i.sendErrorMessage(m, username, err)
				continue
			}

			cursor = accountMedia.Data.User.EdgeOwnerToTimelineMedia.PageInfo.EndCursor
			i.sendUserTimlinePostsId(accountMedia, username, userId)

			if !accountMedia.Data.User.EdgeOwnerToTimelineMedia.PageInfo.HasNextPage {
				isPostsSendingFinished = false
			}
			time.Sleep(time.Millisecond * 100)
		}
		i.nameQReader.CommitMessages(context.Background(), m)
	}
}

func (i *InstaPostsScraper) sendUserInfoPostsId(instagramAccountInfo *instagramAccountInfo, username string, userId string) {
	for _, element := range instagramAccountInfo.Graphql.User.EdgeOwnerToTimelineMedia.Edges {
		fmt.Println("Edges ", username)
		fmt.Println(element.Node.Typename)
		if element.Node.Shortcode != "" {
			fmt.Println("Edges Node1 ", username)
			instagramPost := models.InstagramPost{
				element.Node.ID,
				element.Node.Shortcode,
				userId,
				element.Node.DisplayURL}

			instagramPostJson, err := json.Marshal(instagramPost)

			if err != nil {
				fmt.Println(err)
				break
			}

			err = i.postsQWriter.WriteMessages(context.Background(), kafka.Message{Value: instagramPostJson})
		}
	}
}

func (i *InstaPostsScraper) sendUserTimlinePostsId(accountMedia *InstagramMedia, username string, userId string) {
	for _, element := range accountMedia.Data.User.EdgeOwnerToTimelineMedia.Edges {
		if element.Node.ID != "" {

			instagramPost := models.InstagramPost{
				element.Node.ID,
				element.Node.Shortcode,
				userId,
				element.Node.DisplayURL}

			instagramPostJson, err := json.Marshal(instagramPost)

			if err != nil {
				fmt.Println(err)
				break
			}

			err = i.postsQWriter.WriteMessages(context.Background(), kafka.Message{Value: instagramPostJson})
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}

func (i *InstaPostsScraper) sendErrorMessage(m kafka.Message, username string, err error) {
	errMessage := &models.InstagramScrapeError{
		Name:  username,
		Error: err.Error(),
	}
	serializedErr, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println(err)
	}
	i.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
	i.nameQReader.CommitMessages(context.Background(), m)
}

func (i *InstaPostsScraper) Close() {
	i.Stop()
	i.WaitUntilStopped(time.Second * 3)

	i.nameQReader.Close()
	i.postsQWriter.Close()
	i.errQWriter.Close()
	i.MarkAsClosed()
}
