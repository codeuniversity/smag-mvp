package main

import (
	scraper "github.com/codeuniversity/smag-mvp/insta_posts-scraper"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	awsServiceAddress := utils.MustGetStringFromEnv("AWS_SERVICE_ADDRESS")
	nameReaderConfig, infoWriterConfig, errWriterConfig := kafka.GetInstaPostsScraperConfig()

	s := scraper.New(awsServiceAddress, kafka.NewReader(nameReaderConfig), kafka.NewWriter(infoWriterConfig), kafka.NewWriter(errWriterConfig))

	service.CloseOnSignal(s)
	waitUntilClosed := s.Start()

	waitUntilClosed()
}
