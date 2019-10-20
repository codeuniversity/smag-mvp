package main

import (
	scraper "github.com/codeuniversity/smag-mvp/insta_comments-scraper"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	readerConfig, infoWriterConfig, errWriterConfig := kafka.GetScraperConfig()

	s := scraper.New(kafka.NewReader(readerConfig), kafka.NewWriter(infoWriterConfig), kafka.NewWriter(errWriterConfig))

	service.CloseOnSignal(s)
	waitUntilClosed := s.Start()

	waitUntilClosed()
}
