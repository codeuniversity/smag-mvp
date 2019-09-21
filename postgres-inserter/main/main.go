package main

import (
	"os"

	inserter "github.com/alexmorten/instascraper/postgres-inserter"
	"github.com/alexmorten/instascraper/service"
)

func main() {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "127.0.0.1:9092"
	}
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "127.0.0.1"
	}
	i := inserter.New(kafkaAddress, postgresHost)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
