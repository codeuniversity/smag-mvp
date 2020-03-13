package main

import (
	"github.com/codeuniversity/smag-mvp/config"
	downloader "github.com/codeuniversity/smag-mvp/insta/pics-downloader"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	jobsTopic := utils.MustGetStringFromEnv("KAFKA_PICTURE_DOWNLOADS_TOPIC")
	qReader := kafka.NewReader(kafka.NewReaderConfig(kafkaAddress, groupID, jobsTopic))

	s3Config := config.GetS3Config()
	postgresConfig := config.GetPostgresConfig()

	i := downloader.New(qReader, s3Config, postgresConfig)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}
