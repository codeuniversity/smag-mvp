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
	nameQReader *kafka.Reader
	infoQWriter *kafka.Writer
	errQWriter  *kafka.Writer
}

// New returns an initilized scraper
func New(kafkaAddress string) *Scraper {
	s := &Scraper{}
	s.nameQReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "user_follow_graph_scraper",
		Topic:          "user_names",
		CommitInterval: time.Second,
	})
	s.infoQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "user_follow_infos",
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	s.errQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    "user_scrape_errors",
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})
	return s
}

// Run the scraper
func (s *Scraper) Run() {
	defer s.close()
	fmt.Println("starting scraper")
	for {
		m, err := s.qReader.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}

		userName := string(m.Value)
		followInfo, err := ScrapeUserFollowGraph(userName)
		if err != nil {
			fmt.Println(err)
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

	followerInfo, err := getUserInfoIn(fmt.Sprintf("http://picdeer.com/%s/followers", userName))
	if err != nil {
		return nil, err
	}
	u.Followers = followerInfo.listedUserNames
	followingsInfo, err := getUserInfoIn(fmt.Sprintf("http://picdeer.com/%s/followings", userName))
	if err != nil {
		return nil, err
	}
	u.RealName = followingsInfo.realName
	u.AvatarURL = followingsInfo.avatarURL
	u.Bio = followingsInfo.bio
	u.Followings = followingsInfo.listedUserNames
	u.CrawlTs = int(time.Now().Unix())
	return u, nil
}

type scrapedInfo struct {
	listedUserNames []string
	avatarURL       string
	realName        string
	bio             string
}

func getUserInfoIn(url string) (info *scrapedInfo, err error) {
	info = new(scrapedInfo)

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("p.grid-user-identifier-1 a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			name := parts[len(parts)-1]
			info.listedUserNames = append(info.listedUserNames, name)
		}
	})

	c.OnHTML(".profile-header img.p-avatar", func(e *colly.HTMLElement) {
		info.avatarURL = e.Attr("src")
	})

	c.OnHTML(".profile-header h1.p-h1 a", func(e *colly.HTMLElement) {
		info.realName = e.Text
	})

	c.OnHTML(".profile-header p.p-bio", func(e *colly.HTMLElement) {
		info.bio = e.Text
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err = c.Visit(url)
	if err != nil {
		return nil, err
	}

	return
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
