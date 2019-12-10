package main

import (
	neo4j_dump2 "github.com/codeuniversity/smag-mvp/neo4j-inserter/create-import-json"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")

	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	rTopic := utils.MustGetStringFromEnv("KAFKA_INFO_TOPIC")
	kafkaChunk := utils.GetNumberFromEnvWithDefault("KAFKA_MESSAGE_CHUNK", 10)

	importDump := neo4j_dump2.New(kafkaAddress, rTopic, groupID, kafkaChunk)

	importDump.Run()
}
