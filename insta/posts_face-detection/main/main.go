package main

import (
	detection "github.com/codeuniversity/smag-mvp/insta/posts_face-detection"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	jobsReadTopic := utils.MustGetStringFromEnv("KAFKA_PICTURE_FACE_RECON_TOPIC")
	jobsWriteTopic := utils.MustGetStringFromEnv("KAFKA_PICTURE_FACE_RECONED_TOPIC")
	qReader := kafka.NewReader(kafka.NewReaderConfig(kafkaAddress, groupID, jobsReadTopic))
	qWriter := kafka.NewWriter(kafka.NewWriterConfig(kafkaAddress, jobsWriteTopic, true))

	config := detection.Config{
		S3BucketName:      utils.GetStringFromEnvWithDefault("S3_BUCKET_NAME", "insta_pics"),
		S3Region:          utils.GetStringFromEnvWithDefault("S3_REGION", "eu-west-1"),
		S3Endpoint:        utils.GetStringFromEnvWithDefault("S3_ENDOINT", "127.0.0.1:9000"),
		S3AccessKeyID:     utils.MustGetStringFromEnv("S3_ACCESS_KEY_ID"),
		S3SecretAccessKey: utils.MustGetStringFromEnv("S3_SECRET_ACCESS_KEY"),
		S3UseSSL:          utils.GetBoolFromEnvWithDefault("S3_USE_SSL", true),
	}

	d := detection.New(qReader, qWriter, config)

	service.CloseOnSignal(d)
	waitUntilDone := d.Start()

	waitUntilDone()
}
