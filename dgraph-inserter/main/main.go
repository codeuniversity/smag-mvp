package main

import (
	inserter "github.com/codeuniversity/smag-mvp/dgraph-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	dgraphAddress := utils.GetStringFromEnvWithDefault("DGRPAH_ADDRESS", "127.0.0.1:9080")

	i := inserter.New(kafkaAddress, dgraphAddress)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
