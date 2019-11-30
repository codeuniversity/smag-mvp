package main

import (
	server "github.com/codeuniversity/smag-mvp/api/grpcserver"
	"github.com/codeuniversity/smag-mvp/config"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/utils"
	kgo "github.com/segmentio/kafka-go"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "")
	namesTopic := utils.GetStringFromEnvWithDefault("KAFKA_NAME_TOPIC", "")
	grpcPort := utils.GetStringFromEnvWithDefault("GRPC_PORT", "10000")
	uploadBucket := utils.MustGetStringFromEnv("S3_UPLOAD_BUCKET_NAME")
	esHosts := utils.GetMultipleStringsFromEnvWithDefault("ES_HOSTS", []string{"http://localhost:9200"})
	recognitionServiceAddress := utils.MustGetStringFromEnv("RECOGNITION_SERVICE_ADDRESS")
	s3Config := config.GetS3Config()
	postgresConfig := config.GetPostgresConfig()

	var writer *kgo.Writer
	if kafkaAddress != "" && namesTopic != "" {
		writer = kafka.NewWriter(kafka.NewWriterConfig(kafkaAddress, namesTopic, false))
	}

	s := server.NewGrpcServer(grpcPort, writer, s3Config, uploadBucket, postgresConfig, esHosts, recognitionServiceAddress)

	s.Listen()
}
