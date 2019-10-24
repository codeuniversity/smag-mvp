package main

import (
	filter "github.com/codeuniversity/smag-mvp/insta_user_names-filter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.users")
	namesTopic := utils.GetStringFromEnvWithDefault("KAFKA_NAME_TOPIC", "user_names")

	f := filter.New(kafkaAddress, groupID, changesTopic, namesTopic)

	service.CloseOnSignal(f)
	waitUntilClosed := f.Start()

	waitUntilClosed()
}
