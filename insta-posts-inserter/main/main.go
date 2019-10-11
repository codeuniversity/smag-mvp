package main

import (
	"github.com/codeuniversity/smag-mvp/insta-posts-inserter"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	postgresHost := utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1")
	postgresPassword := utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", "")

	qReaderConfig, qWriterConfig, _ := kafka.GetInserterConfig()

	i := insta_posts_inserter.New(postgresHost, postgresPassword, kafka.NewReader(qReaderConfig), kafka.NewWriter(qWriterConfig))

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
