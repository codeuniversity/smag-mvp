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
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cli/main/main.go <instagram|twitter> <username>")
		return
	}

	platformArg := os.Args[1]
	userNameArg := os.Args[2]

	var topic string
	switch platformArg {
	case "instagram":
		topic = "user_names"
		break
	case "twitter":
		topic = "twitter-user_names"
		break
	default:
		fmt.Printf("Invalid platform option: %s\n", platformArg)
		return
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()
	t, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := w.WriteMessages(t, kafka.Message{
		Value: []byte(userNameArg),
	})
	utils.PanicIfErr(err)
}
