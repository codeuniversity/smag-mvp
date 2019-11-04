package main

import (
	"github.com/codeuniversity/smag-mvp/aws_service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	grpcPort := utils.GetStringFromEnvWithDefault("GRPC_PORT", "9900")
	s := aws_service.New(grpcPort)
	s.Listen()
}
