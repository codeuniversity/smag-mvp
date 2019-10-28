package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/segmentio/kafka-go"
)

func main() {
	var platformArg string
	var userNameArg string

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	instagramTopic := utils.GetStringFromEnvWithDefault("INSTAGRAM_TOPIC", "user_names")
	twitterTopic := utils.GetStringFromEnvWithDefault("Twitter_TOPIC", "twitter-user_names")

	if len(os.Args) == 3 {
		fmt.Println("getting arguments from parameters")
		platformArg = os.Args[1]
		userNameArg = os.Args[2]
	} else {
		fmt.Println("getting arguments from env variables")
		platformArg = utils.MustGetStringFromEnv("CLI_PLATFORM")
		userNameArg = utils.MustGetStringFromEnv("CLI_USER")
	}

	var topic string
	switch platformArg {
	case "instagram":
		topic = instagramTopic
		break
	case "twitter":
		topic = twitterTopic
		break
	default:
		fmt.Printf("Invalid platform option: %s\n", platformArg)
		return
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()
	t, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := w.WriteMessages(t, kafka.Message{
		Value: []byte(userNameArg),
	})
	utils.PanicIfNotNil(err)
}
