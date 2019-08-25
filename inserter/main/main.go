package main

import (
	"os"

	"github.com/alexmorten/instascraper/inserter"
	"github.com/alexmorten/instascraper/utils"
)

func main() {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "127.0.0.1:9092"
	}
	dgraphAddress := os.Getenv("DGRAPH_ADDRESS")
	if dgraphAddress == "" {
		dgraphAddress = "127.0.0.1:9080"
	}
	i := inserter.New(kafkaAddress, dgraphAddress)

	utils.CloseOnSignal(i.Close)
	go i.Run()

	i.WaitUntilClosed()
}
