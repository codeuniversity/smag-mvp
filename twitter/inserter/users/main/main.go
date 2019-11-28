package main

import (
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	inserter "github.com/codeuniversity/smag-mvp/twitter/inserter/users"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	var i *inserter.Inserter

	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")
	postgresPassword := utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", "")
	postgresDBNAme := utils.GetStringFromEnvWithDefault("POSTGRES_DB_NAME", "twitter_scraper")

	qReaderConfig := kafka.GetInserterConfig()

	i = inserter.New(
		postgresHost,
		postgresPassword,
		postgresDBNAme,
		kafka.NewReader(qReaderConfig),
	)

	service.CloseOnSignal(i)
	waitUntilDone := i.Start()
	waitUntilDone()
}
