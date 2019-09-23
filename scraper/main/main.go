package main

import (
	"os"

	"github.com/codeuniversity/smag-mvp/scraper"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "localhost:9092"
	}
	s := scraper.New(kafkaAddress)
	service.CloseOnSignal(s)
	go s.Run()

	s.WaitUntilClosed()
}
