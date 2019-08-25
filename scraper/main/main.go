package main

import (
	"os"

	"github.com/alexmorten/instascraper/scraper"
	"github.com/alexmorten/instascraper/utils"
)

func main() {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "localhost:9092"
	}
	s := scraper.New(kafkaAddress)
	utils.CloseOnSignal(s.Close)
	go s.Run()

	s.WaitUntilClosed()
}
