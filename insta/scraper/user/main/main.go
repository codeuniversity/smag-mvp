package main

import (
	scraper "github.com/codeuniversity/smag-mvp/insta/scraper/user"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	nameReaderConfig, infoWriterConfig, errWriterConfig := kafka.GetScraperConfig()

	s := scraper.New(
		kafka.NewReader(nameReaderConfig),
		kafka.NewWriter(infoWriterConfig),
		kafka.NewWriter(errWriterConfig),
	)
	service.CloseOnSignal(s)
	waitUntilDone := s.Start()

	waitUntilDone()
}
