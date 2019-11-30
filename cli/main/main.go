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
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	instagramTopic := utils.GetStringFromEnvWithDefault("KAFKA_INSTAGRAM_TOPIC", "user_names")
	twitterTopic := utils.GetStringFromEnvWithDefault("KAFKA_TWITTER_TOPIC", "twitter.scraped.user_names")

	if len(os.Args) < 3 {
		panic("Invalid argumemts. Usage: cli <instagram|twitter> <username>")
	}

	platformArg := os.Args[1]
	userNameArg := os.Args[2]

	var topic string
	switch platformArg {
	case "instagram":
		topic = instagramTopic
		break
	case "twitter":
		topic = twitterTopic
		break
	default:
		panic(fmt.Sprintf("Invalid platform option: %s\n", platformArg))
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
