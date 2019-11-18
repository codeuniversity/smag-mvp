package main

import (
	inserter "github.com/codeuniversity/smag-mvp/insta_neo4j_user_names-inserter"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	neo4jHost := utils.GetStringFromEnvWithDefault("NEO4J_HOST", "127.0.0.1")
	neo4jUsername := utils.GetStringFromEnvWithDefault("NEO4J_USERNAME", "neo4j")
	neo4jPassword := utils.GetStringFromEnvWithDefault("NEO4J_PASSWORD", "neo4j")

	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.users")
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	qReaderConfig := kafka.NewReaderConfig(kafkaAddress, groupID, changesTopic)

	i := inserter.New(neo4jHost, neo4jUsername, neo4jPassword, kafka.NewReader(qReaderConfig))

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}
