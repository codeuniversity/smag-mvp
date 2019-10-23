package main

import (
	inserter "github.com/codeuniversity/smag-mvp/insta_postgres-inserter"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	var i *inserter.Inserter

	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")
	postgresPassword := utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", "")

	qReaderConfig, qWriterConfig, isUserDiscovery := kafka.GetUserDiscoveryInserterConfig()

	if isUserDiscovery {
		i = inserter.New(
			postgresHost,
			postgresPassword,
			kafka.NewReader(qReaderConfig),
			kafka.NewWriter(qWriterConfig),
		)
	} else {
		i = inserter.New(
			postgresHost,
			postgresPassword,
			kafka.NewReader(qReaderConfig),
			nil,
		)
	}

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}
