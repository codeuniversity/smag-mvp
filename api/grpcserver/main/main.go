package main

import (
	server "github.com/codeuniversity/smag-mvp/api/grpcserver"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	grpcPort := utils.GetStringFromEnvWithDefault("GRPC_PORT", "10000")

	config := models.Config{
		S3BucketName:      utils.GetStringFromEnvWithDefault("S3_BUCKET_NAME", "insta_pics"),
		S3Region:          utils.GetStringFromEnvWithDefault("S3_REGION", "eu-west-1"),
		S3Endpoint:        utils.GetStringFromEnvWithDefault("S3_ENDOINT", "127.0.0.1:9000"),
		S3AccessKeyID:     utils.MustGetStringFromEnv("S3_ACCESS_KEY_ID"),
		S3SecretAccessKey: utils.MustGetStringFromEnv("S3_SECRET_ACCESS_KEY"),
		S3UseSSL:          utils.GetBoolFromEnvWithDefault("S3_USE_SSL", true),

		PostgresHost:     utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1"),
		PostgresPassword: utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", ""),
	}

	s := server.NewGrpcServer(grpcPort, config)

	s.Listen()

}
