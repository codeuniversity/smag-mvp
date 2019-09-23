package main

import (
	"os"

	"github.com/codeuniversity/smag-mvp/dgraph-inserter"
	"github.com/codeuniversity/smag-mvp/service"
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

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
