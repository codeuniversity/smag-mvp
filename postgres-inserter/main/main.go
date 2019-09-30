package main

import (
	inserter "github.com/codeuniversity/smag-mvp/postgres-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")
	postgresPassword := utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", "")

	i := inserter.New(kafkaAddress, postgresHost, postgresPassword)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
