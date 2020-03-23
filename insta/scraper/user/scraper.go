package insta_user_scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	http_header_generator "github.com/codeuniversity/smag-mvp/http_header-generator"
	"github.com/codeuniversity/smag-mvp/insta/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"

	"github.com/gocolly/colly"
	"github.com/segmentio/kafka-go"
)

// Scraper represents the scraper containing all clients it uses
type Scraper struct {
	*worker.Worker

	nameQReader *kafka.Reader
	infoQWriter *kafka.Writer
	errQWriter  *kafka.Writer
	*http_header_generator.HTTPHeaderGenerator
}

// New returns an initilized scraper
func New(nameQReader *kafka.Reader, infoQWriter *kafka.Writer, errQWriter *kafka.Writer) *Scraper {
	s := &Scraper{}
	s.nameQReader = nameQReader
	s.infoQWriter = infoQWriter
	s.errQWriter = errQWriter
	s.HTTPHeaderGenerator = http_header_generator.New()

	s.Worker = worker.Builder{}.WithName("insta_user-scraper").
		WithWorkStep(s.runStep).
		AddShutdownHook("nameQReader", nameQReader.Close).
		AddShutdownHook("infoQWriter", infoQWriter.Close).
		AddShutdownHook("errQWriter", errQWriter.Close).
		MustBuild()
	return s
}

// runStep the scraper
func (s *Scraper) runStep() error {
	log.Println("fetching")
	m, err := s.nameQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	userName := string(m.Value)
	followInfo, err := s.scrapeUserFollowGraph(userName)
	if err != nil {
		log.Println(err)
		errMessage := &models.ScrapeError{
			Name:  userName,
			Error: err.Error(),
		}
		serializedErr, serializationErr := json.Marshal(errMessage)
		if serializationErr != nil {
			return serializationErr
		}
		err := s.errQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedErr})
		if err != nil {
			return err
		}

		return s.nameQReader.CommitMessages(context.Background(), m)
	}
	serializedFollowInfo, err := json.Marshal(followInfo)
	if err != nil {
		return err
	}
	err = s.infoQWriter.WriteMessages(context.Background(), kafka.Message{Value: serializedFollowInfo})
	if err != nil {
		return err
	}
	s.nameQReader.CommitMessages(context.Background(), m)

	return nil
}

//scrapeUserFollowGraph returns the follow information for a userName
func (s *Scraper) scrapeUserFollowGraph(userName string) (*models.UserFollowInfo, error) {
	u := &models.UserFollowInfo{UserName: userName}

	err := utils.WithRetries(5, func() error {
		followingsInfo, err := s.getUserInfoIn(fmt.Sprintf("http://picdeer.com/%s/followings", userName))
		if err != nil {
			return err
		}

		u.RealName = followingsInfo.realName
		u.AvatarURL = followingsInfo.avatarURL
		u.Bio = followingsInfo.bio
		u.Followings = followingsInfo.listedUserNames
		u.CrawlTs = int(time.Now().Unix())
		return nil
	})

	if err != nil {
		return nil, err
	}

	return u, nil
}

type scrapedInfo struct {
	listedUserNames []string
	avatarURL       string
	realName        string
	bio             string
}

func (s *Scraper) getUserInfoIn(url string) (info *scrapedInfo, err error) {
	info = new(scrapedInfo)

	c := colly.NewCollector()
	c.UserAgent = s.GetRandomUserAgent()
	c.OnRequest(func(r *colly.Request) {
		s.AddHeaders(r.Headers)
	})

	c.OnError(func(c *colly.Response, err error) {
		log.Println("Something went wrong:", err, " code: ", c.StatusCode)
	})

	c.OnResponse(func(r *colly.Response) {
		// log.Println("Visited", r.Request.URL)
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
		log.Println("Finished", r.Request.URL)
	})

	err = c.Visit(url)
	if err != nil {
		return nil, err
	}

	return
}
