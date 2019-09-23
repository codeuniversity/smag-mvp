package main

import (
	"os"

	inserter "github.com/codeuniversity/smag-mvp/neo4j-inserter"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "127.0.0.1:9092"
	}
	neo4jAddress := os.Getenv("NEO4J_ADDRESS")
	neo4jUsername := os.Getenv("NEO4J_USERNAME")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")

	if neo4jAddress == "" {
		neo4jAddress = "127.0.0.1:7687"
	}

	if neo4jUsername == "" {
		neo4jUsername = "neo4j"
	}

	if neo4jPassword == "" {
		neo4jPassword = "123456"
	}
	i := inserter.New(kafkaAddress, neo4jAddress, neo4jUsername, neo4jPassword)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
