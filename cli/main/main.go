package main

import (
	"context"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a username as the first argument")
		return
	}
	userNameArg := os.Args[1]
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"172.31.32.93:9092"},
		Topic:    "user_names",
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(userNameArg),
	})
}
