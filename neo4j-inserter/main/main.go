package main

import (
	"github.com/codeuniversity/smag-mvp/kafka"
	inserter "github.com/codeuniversity/smag-mvp/neo4j-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	var i *inserter.Inserter

	neo4jAddress := utils.GetStringFromEnvWithDefault("NEO4J_ADDRESS", "127.0.0.1:7687")
	neo4jUsername := utils.GetStringFromEnvWithDefault("NEO4J_USERNAME", "neo4j")
	neo4jPassword := utils.GetStringFromEnvWithDefault("NEO4J_PASSWORD", "123456")

	qReaderConfig, qWriterConfig, isUserDiscovery := kafka.GetInserterConfig()

	if isUserDiscovery {
		i = inserter.New(
			neo4jAddress,
			neo4jUsername,
			neo4jPassword,
			kafka.NewReader(qReaderConfig),
			kafka.NewWriter(qWriterConfig),
		)
	} else {
		i = inserter.New(
			neo4jAddress,
			neo4jUsername,
			neo4jPassword,
			kafka.NewReader(qReaderConfig),
			nil,
		)
	}

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
