package main

import (
	downloader "github.com/codeuniversity/smag-mvp/insta_pics-downloader"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	jobsTopic := utils.MustGetStringFromEnv("KAFKA_PICTURE_DOWNLOADS_TOPIC")
	qReader := kafka.NewReader(kafka.NewReaderConfig(kafkaAddress, groupID, jobsTopic))

	config := downloader.Config{
		S3BucketName:      utils.GetStringFromEnvWithDefault("S3_BUCKET_NAME", "insta_pics"),
		S3Region:          utils.GetStringFromEnvWithDefault("S3_REGION", "eu-west-1"),
		S3Endpoint:        utils.GetStringFromEnvWithDefault("S3_ENDOINT", "127.0.0.1:9000"),
		S3AccessKeyID:     utils.MustGetStringFromEnv("S3_ACCESS_KEY_ID"),
		S3SecretAccessKey: utils.MustGetStringFromEnv("S3_SECRET_ACCESS_KEY"),
		S3UseSSL:          utils.GetBoolFromEnvWithDefault("S3_USE_SSL", true),

		PostgresHost:     utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1"),
		PostgresPassword: utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", ""),
	}

	i := downloader.New(qReader, config)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}
