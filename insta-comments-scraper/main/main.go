package main

import (
	insta_comments_scraper "github.com/codeuniversity/smag-mvp/insta-comments-scraper"
	"github.com/codeuniversity/smag-mvp/service"
	"os"
)

func main() {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "52.58.171.160:9092"
	}
	s := insta_comments_scraper.New(kafkaAddress)
	service.CloseOnSignal(s)
	go s.Run()

	s.WaitUntilClosed()
}
