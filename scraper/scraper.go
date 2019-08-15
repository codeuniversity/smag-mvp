package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alexmorten/instascraper/models"
	"github.com/gocolly/colly"
	"github.com/segmentio/kafka-go"
)

// Scraper represents the scraper containing all clients it uses
type Scraper struct {
	qReader *kafka.Reader
	qWriter *kafka.Writer
}

// New returns an initilized scraper
func New() *Scraper {
	s := &Scraper{}
	s.qReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		GroupID:        "user_follow_graph_scraper",
		Topic:          "user_names",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
	s.qWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "user_follow_infos",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	return s
}

// Run the scraper
func (s *Scraper) Run() {
	defer s.close()

	for {
		m, err := s.qReader.FetchMessage(context.Background())
		if err != nil {
			break
		}

		userName := string(m.Value)
		followInfo, err := ScrapeUserFollowGraph(userName)
		if err != nil {
			break
		}
		serializedFollowInfo, err := json.Marshal(followInfo)
		if err != nil {
			fmt.Println(err)
			break
		}
		err = s.qWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedFollowInfo})
		if err != nil {
			fmt.Println(err)
			break
		}
		s.qReader.CommitMessages(context.Background(), m)
	}
}

func (s *Scraper) close() {
	s.qReader.Close()
	s.qWriter.Close()
}

//ScrapeUserFollowGraph returns the follow information for a userName
func ScrapeUserFollowGraph(userName string) (*models.UserFollowInfo, error) {
	u := &models.UserFollowInfo{UserName: userName}

	followers, err := getUserNamesIn(fmt.Sprintf("http://picdeer.com/%s/followers", userName))
	if err != nil {
		return nil, err
	}
	u.Followers = followers
	followings, err := getUserNamesIn(fmt.Sprintf("http://picdeer.com/%s/followings", userName))
	if err != nil {
		return nil, err
	}
	u.Followings = followings
	return u, nil
}

func getUserNamesIn(url string) ([]string, error) {
	userNames := []string{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		// log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("p.grid-user-identifier-1 a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			name := parts[len(parts)-1]
			userNames = append(userNames, name)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return userNames, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
