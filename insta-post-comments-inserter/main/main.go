package main

import (
	"github.com/codeuniversity/smag-mvp/insta-post-comments-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "52.58.171.160:9092")

	s := insta_post_comments_inserter.New(kafkaAddress, postgresHost)
	service.CloseOnSignal(s)
	go s.Run()

	s.WaitUntilClosed()
}
