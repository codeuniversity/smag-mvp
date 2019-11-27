package main

import (
	server "github.com/codeuniversity/smag-mvp/api/grpcserver"
	"github.com/codeuniversity/smag-mvp/config"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	grpcPort := utils.GetStringFromEnvWithDefault("GRPC_PORT", "10000")
	uploadBucket := utils.MustGetStringFromEnv("S3_UPLOAD_BUCKET_NAME")
	esHosts := utils.GetMultipleStringsFromEnvWithDefault("ES_HOSTS", []string{"http://localhost:9200"})
	recognitionServiceAddress := utils.MustGetStringFromEnv("RECOGNITION_SERVICE_ADDRESS")
	s3Config := config.GetS3Config()
	postgresConfig := config.GetPostgresConfig()

	s := server.NewGrpcServer(grpcPort, s3Config, uploadBucket, postgresConfig, esHosts, recognitionServiceAddress)

	s.Listen()
}
