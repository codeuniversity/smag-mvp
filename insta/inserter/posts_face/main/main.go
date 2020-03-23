package main

import (
	inserter "github.com/codeuniversity/smag-mvp/insta/inserter/posts_face"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	var i *inserter.Inserter

	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")
	postgresPassword := utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", "")

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	jobsReadTopic := utils.GetStringFromEnvWithDefault("KAFKA_FACE_DETECTION_RESULTS_TOPIC", "insta_posts_detected_faces")
	qReader := kafka.NewReader(kafka.NewReaderConfig(kafkaAddress, groupID, jobsReadTopic))

	i = inserter.New(
		postgresHost,
		postgresPassword,
		qReader,
	)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}
