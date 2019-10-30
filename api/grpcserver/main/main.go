package main

import (
	server "github.com/codeuniversity/smag-mvp/api/grpcserver"
	"github.com/codeuniversity/smag-mvp/config"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	grpcPort := utils.GetStringFromEnvWithDefault("GRPC_PORT", "10000")

	s3Config := config.GetS3Config()
	postgresConfig := config.GetPostgresConfig()

	s := server.NewGrpcServer(grpcPort, s3Config, postgresConfig)

	s.Listen()

}
