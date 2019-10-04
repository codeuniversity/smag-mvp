package instagramScraper

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

type InstagramScraper struct {
	nameQReader      *kafka.Reader
	userPostsQWriter *kafka.Writer
	errQWriter       *kafka.Writer
	*service.Executor
	kafkaAddress string
	httpClient   *httpClient.HttpClient
}

// New returns an initilized scraper
func New(kafkaAddress string) *InstagramScraper {
	i := &InstagramScraper{}
	i.nameQReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "user_instagram_scraper_posts",
		Topic:          "user_names",
		CommitInterval: time.Minute * 40,
	})
	i.userPostsQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "user_post",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	i.errQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "user_instagram-scrape_errors",
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})
	i.Executor = service.New()
	i.kafkaAddress = kafkaAddress
	i.httpClient = httpClient.New(2, i.kafkaAddress)
	return i
}

func (i *InstagramScraper) accountInfo(username string) (*models.InstagramAccountInfo, error) {
	var instagramAccountInfo *models.InstagramAccountInfo

	err := i.httpClient.WithRetries(2, func() error {
		accountInfo, err := i.httpClient.ScrapeAccountInfo(username)
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

func (i *InstagramScraper) accountPosts(userId string, cursor string) (*models.InstagramMedia, error) {
	var instagramAccountMedia *models.InstagramMedia

	err := i.httpClient.WithRetries(2, func() error {
		accountInfo, err := i.httpClient.ScrapeProfileMedia(userId, cursor)
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

func (i *InstagramScraper) Run() {
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

		for _, element := range instagramAccountInfo.Graphql.User.EdgeOwnerToTimelineMedia.Edges {
			fmt.Println("Edges ", username)
			fmt.Println(element.Node.Typename)
			if element.Node.Shortcode != "" {
				fmt.Println("Edges Node1 ", username)
				instagramPost := models.InstagramPost{
					element.Node.ID,
					element.Node.Shortcode,
					userId}

				instagramPostJson, err := json.Marshal(instagramPost)

				if err != nil {
					fmt.Println(err)
					break
				}

				err = i.userPostsQWriter.WriteMessages(context.Background(), kafka.Message{Value: instagramPostJson})
			} else if element.Node.Shortcode != "" {
				fmt.Println("Edges Node2 ", username)
				instagramPost := models.InstagramPost{
					element.Node.ID,
					element.Node.Shortcode,
					userId}

				instagramPostJson, err := json.Marshal(instagramPost)

				if err != nil {
					fmt.Println(err)
					break
				}

				err = i.userPostsQWriter.WriteMessages(context.Background(), kafka.Message{Value: instagramPostJson})
			}
		}

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

			for _, element := range accountMedia.Data.User.EdgeOwnerToTimelineMedia.Edges {
				if element.Node.ID != "" {

					instagramPost := models.InstagramPost{
						element.Node.ID,
						element.Node.Shortcode,
						userId}

					instagramPostJson, err := json.Marshal(instagramPost)

					if err != nil {
						fmt.Println(err)
						break
					}

					err = i.userPostsQWriter.WriteMessages(context.Background(), kafka.Message{Value: instagramPostJson})
					if err != nil {
						fmt.Println(err)
						break
					}
				}
			}

			cursor = accountMedia.Data.User.EdgeOwnerToTimelineMedia.PageInfo.EndCursor
			if !accountMedia.Data.User.EdgeOwnerToTimelineMedia.PageInfo.HasNextPage {
				isPostsSendingFinished = false
			}
			time.Sleep(time.Millisecond * 600)
		}
		i.nameQReader.CommitMessages(context.Background(), m)
	}
}

func (i *InstagramScraper) sendErrorMessage(m kafka.Message, username string, err error) {
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

func (i *InstagramScraper) Close() {
	i.Stop()
	i.WaitUntilStopped(time.Second * 3)

	i.nameQReader.Close()
	i.userPostsQWriter.Close()
	i.errQWriter.Close()
	i.httpClient.Close()
	i.MarkAsClosed()
}
