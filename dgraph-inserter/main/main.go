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
	isUserDiscovery := utils.GetBoolFromEnvWithDefault("USER_DISCOVERY", false)

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	dgraphAddress := utils.GetStringFromEnvWithDefault("DGRPAH_ADDRESS", "127.0.0.1:9080")

	kafkaConfig := utils.NewKafkaConsumerConfig(groupID, rTopic, wTopic, isUserDiscovery)
	i := inserter.New(kafkaAddress, dgraphAddress, kafkaConfig)

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
