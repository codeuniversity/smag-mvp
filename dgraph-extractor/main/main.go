package main

import (
	"os"
	"strconv"

	extractor "github.com/alexmorten/instascraper/dgraph-extractor"
	"github.com/alexmorten/instascraper/service"
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

	startID := os.Getenv("START_ID")
	if startID == "" {
		startID = "1"
	}
	startIDInt, err := strconv.ParseInt(startID, 10, 64)
	if err != nil {
		panic(err)
	}
	i := extractor.New(kafkaAddress, dgraphAddress, int(startIDInt))

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
