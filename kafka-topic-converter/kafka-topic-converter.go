package kafka_topic_converter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/segmentio/kafka-go"
)

type Converter struct {
	kafkaTopicIn  *kafka.Reader
	kafkaTopicOut *kafka.Writer

	*service.Executor
}

func New(topicIn *kafa.Reader, topicOut *kafka.Writer) *Converter {
	c := &Converter{}
	c.kafkaTopicIn = topicIn
	c.kafkaTopicOut = topicOut

	c.Executor = service.New()

	return c
}

func (c *Converter) Run() {

	for c.IsRunning() {
		m, err := c.kafkaTopicIn.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}

		changeStream := &models.DebeziumChangeStream{}
		err = json.Unmarshal(m.Value, changeStream)
		if err != nil {
			fmt.Println(err)
			break
		}
		userName := changeStream.Payload.After.UserName
		c.kafkaTopicIn.CommitMessages(context.Background(), m)

		c.kafkaTopicOut.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(userName),
		})

		fmt.Printf("%s tranfered \n", userName)

	}
}
