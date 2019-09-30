package main

import (
	inserter "github.com/codeuniversity/smag-mvp/dgraph-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	rTopic := utils.MustGetStringFromEnv("KAFKA_RTOPIC")
	wTopic := utils.MustGetStringFromEnv("KAFKA_WTOPIC")

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	dgraphAddress := utils.GetStringFromEnvWithDefault("DGRPAH_ADDRESS", "127.0.0.1:9080")

	i := inserter.New(kafkaAddress, dgraphAddress, groupID, rTopic, wTopic)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
