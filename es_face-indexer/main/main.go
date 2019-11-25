package main

import (
	indexer "github.com/codeuniversity/smag-mvp/es_face-indexer"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
	"strings"
)

func main() {
	esHostsString := utils.GetStringFromEnvWithDefault("ES_HOSTS", "http://localhost:9200")
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.face_data")

	changeQReader := kafka.NewReader(kafka.NewReaderConfig(kafkaAddress, groupID, changesTopic))
	esHosts := strings.Split(esHostsString, ",")
	i := indexer.New(changeQReader, esHosts)

	service.CloseOnSignal(i)
	waitUntilDone := i.Start()

	waitUntilDone()
}
