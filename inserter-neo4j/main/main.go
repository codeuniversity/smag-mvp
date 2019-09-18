package main

import (
	"os"

	inserter "github.com/alexmorten/instascraper/inserter-neo4j"
	"github.com/alexmorten/instascraper/service"
)

func main() {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "127.0.0.1:9092"
	}
	neo4jAddress := os.Getenv("NEO4J_ADDRESS")
	if neo4jAddress == "" {
		neo4jAddress = "bolt://neo4j:123456@127.0.0.1:7687"
	}
	i := inserter.New(kafkaAddress, neo4jAddress)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
