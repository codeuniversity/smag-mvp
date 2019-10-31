package main

import (
	scraper "github.com/codeuniversity/smag-mvp/insta_comments-scraper"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	awsServiceAddress := utils.GetStringFromEnvWithDefault("AWS_SERVICE_ADDRESS", "")
	readerConfig, infoWriterConfig, errWriterConfig := kafka.GetScraperConfig()

	s := scraper.New(awsServiceAddress, kafka.NewReader(readerConfig), kafka.NewWriter(infoWriterConfig), kafka.NewWriter(errWriterConfig))

	service.CloseOnSignal(s)
	waitUntilClosed := s.Start()

	waitUntilClosed()
}
