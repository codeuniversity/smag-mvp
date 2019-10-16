package main

import (
	insta_posts_scraper "github.com/codeuniversity/smag-mvp/insta-posts-scraper"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	nameReaderConfig, infoWriterConfig, errWriterConfig := kafka.GetScraperConfig()

	s := insta_posts_scraper.New(kafka.NewReader(nameReaderConfig), kafka.NewWriter(infoWriterConfig), kafka.NewWriter(errWriterConfig))

	service.CloseOnSignal(s)
	go s.Run()

	s.WaitUntilClosed()
}
