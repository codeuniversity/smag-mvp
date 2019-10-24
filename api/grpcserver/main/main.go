package main

import (
	server "github.com/codeuniversity/smag-mvp/api/grpcserver"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")
	postgresPassword := utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", "")

	s := server.NewGrpcServer(postgresHost, postgresPassword, 10000)

	s.Listen()

}
