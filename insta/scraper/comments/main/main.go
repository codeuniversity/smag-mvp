package main

import (
	scraper "github.com/codeuniversity/smag-mvp/insta/scraper/comments"
	"github.com/codeuniversity/smag-mvp/kafka"
	client "github.com/codeuniversity/smag-mvp/scraper-client"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	awsServiceAddress := utils.GetStringFromEnvWithDefault("AWS_SERVICE_ADDRESS", "")
	commentLimit := utils.GetNumberFromEnvWithDefault("COMMENT_LIMIT", 24)
	readerConfig, infoWriterConfig, errWriterConfig := kafka.GetScraperConfig()

	config := client.GetScraperConfig()
	s := scraper.New(config, awsServiceAddress, kafka.NewReader(readerConfig), kafka.NewWriter(infoWriterConfig), kafka.NewWriter(errWriterConfig), commentLimit)

	service.CloseOnSignal(s)
	waitUntilClosed := s.Start()

	waitUntilClosed()
}
