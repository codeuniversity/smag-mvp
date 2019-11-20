package main

import (
	"strings"

	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.GetStringFromEnvWithDefault("KAFKA_GROUPID", "insta_user-elastic_filter")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.users")

	elasticAddresses := utils.GetStringFromEnvWithDefault("ELASTIC_ADDRESSES", "http://localhost:9200,http://localhost:9201")
	splitAddr := strings.Split(elasticAddresses, ",")

	f := changestream.NewKafkaToElasticFilter(kafkaAddress, groupID, changesTopic, splitAddr)

	service.CloseOnSignal(f)
	waitUntilClosed := f.Start()

	waitUntilClosed()
}
