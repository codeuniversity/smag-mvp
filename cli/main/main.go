package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a username as the first argument")
		return
	}
	userNameArg := os.Args[1]
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "user_names",
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()
	t, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := w.WriteMessages(t, kafka.Message{
		Value: []byte(userNameArg),
	})
	if err != nil {
		panic(err)
	}
}
