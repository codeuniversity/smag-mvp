package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

func main() {

	kafkaAddress := "52.58.171.160:9092"
	renewedAddressQReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddress},
		GroupID:        "instagram_group1",
		Topic:          "renewed_elastic_ip",
		CommitInterval: time.Minute * 10,
	})

	counter := 0
	for true {
		m, err := renewedAddressQReader.FetchMessage(context.Background())

		if err != nil {
			fmt.Println("Err: ", err)
		}

		fmt.Println("Counter: ", counter)
		fmt.Println("MessageTime: ", m.Time)
		counter++
	}
}
