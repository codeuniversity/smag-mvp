package main

import (
	recognition "github.com/codeuniversity/smag-mvp/face-recognition"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	jobsReadTopic := utils.MustGetStringFromEnv("KAFKA_PICTURE_FACE_RECONED_TOPIC")
	jobsWriteTopic := utils.MustGetStringFromEnv("KAFKA_FACE_DETECTION_RESULTS_TOPIC")
	faceRecognizerAddress := utils.MustGetStringFromEnv("FACE_RECOGNIZER_ADDRESS")
	pictureBucketName := utils.MustGetStringFromEnv("S3_PICTURE_BUCKET_NAME")
	imgProxyAddress := utils.MustGetStringFromEnv("IMGPROXY_ADDRESS")
	imgProxyKey := utils.MustGetStringFromEnv("IMGPROXY_KEY")
	imgProxySalt := utils.MustGetStringFromEnv("IMGPROXY_SALT")
	qReader := kafka.NewReader(kafka.NewReaderConfig(kafkaAddress, groupID, jobsReadTopic))
	qWriter := kafka.NewWriter(kafka.NewWriterConfig(kafkaAddress, jobsWriteTopic, true))

	r := recognition.New(qReader, qWriter, faceRecognizerAddress, pictureBucketName, imgProxyAddress, imgProxyKey, imgProxySalt)

	service.CloseOnSignal(r)
	waitUntilDone := r.Start()

	waitUntilDone()
}
