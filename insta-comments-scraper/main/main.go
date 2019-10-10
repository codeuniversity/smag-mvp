package main

import (
	insta_comments_scraper "github.com/codeuniversity/smag-mvp/insta-comments-scraper"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "52.58.171.160:9092")
	s := insta_comments_scraper.New(kafkaAddress)
	service.CloseOnSignal(s)
	go s.Run()

	s.WaitUntilClosed()
}
