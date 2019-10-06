package main

import (
	inserter "github.com/codeuniversity/smag-mvp/postgres-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "172.31.32.93:9092")
	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")

	i := inserter.New(kafkaAddress, postgresHost)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
