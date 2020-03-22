package main

import (
	neo4j_import "github.com/codeuniversity/smag-mvp/neo4j/create-import-user-json"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")

	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	rTopic := utils.MustGetStringFromEnv("KAFKA_CHANGE_TOPIC")
	kafkaChunk := utils.GetNumberFromEnvWithDefault("KAFKA_MESSAGE_CHUNK", 10)

	i := neo4j_import.New(kafkaAddress, rTopic, groupID, kafkaChunk)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}
