package main

import (
	insta_comments_scraper "github.com/codeuniversity/smag-mvp/insta-comments-scraper"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	readerConfig, infoWriterConfig, errWriterConfig := kafka.GetScraperConfig()

	s := insta_comments_scraper.New(kafka.NewReader(readerConfig), kafka.NewWriter(infoWriterConfig), kafka.NewWriter(errWriterConfig))

	service.CloseOnSignal(s)
	go s.Run()

	s.WaitUntilClosed()
}
