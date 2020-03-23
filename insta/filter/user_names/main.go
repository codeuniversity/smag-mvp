package main

import (
	"encoding/json"

	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"

	"github.com/segmentio/kafka-go"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.users")
	namesTopic := utils.GetStringFromEnvWithDefault("KAFKA_NAME_TOPIC", "user_names")

	f := changestream.NewFilter(kafkaAddress, groupID, changesTopic, namesTopic, filterChange)

	service.CloseOnSignal(f)
	waitUntilClosed := f.Start()

	waitUntilClosed()
}

type user struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
}

func filterChange(m *changestream.ChangeMessage) ([]kafka.Message, error) {
	if m.Payload.Op != "c" {
		return nil, nil
	}

	u := &user{}
	err := json.Unmarshal(m.Payload.After, u)
	if err != nil {
		return nil, err
	}

	return []kafka.Message{{Value: []byte(u.UserName)}}, nil
}
