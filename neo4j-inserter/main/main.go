package main

import (
	inserter "github.com/codeuniversity/smag-mvp/neo4j-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	neo4jAddress := utils.GetStringFromEnvWithDefault("NEO4J_ADDRESS", "127.0.0.1:7687")
	neo4jUsername := utils.GetStringFromEnvWithDefault("NEO4J_USERNAME", "neo4j")
	neo4jPassword := utils.GetStringFromEnvWithDefault("NEO4J_PASSWORD", "123456")

	i := inserter.New(kafkaAddress, neo4jAddress, neo4jUsername, neo4jPassword)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
