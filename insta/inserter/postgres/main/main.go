package main

import (
	inserter "github.com/codeuniversity/smag-mvp/insta/inserter/postgres"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	var i *inserter.Inserter

	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")
	postgresPassword := utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", "")

	qReaderConfig := kafka.GetInserterConfig()

	i = inserter.New(
		postgresHost,
		postgresPassword,
		kafka.NewReader(qReaderConfig),
	)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}
