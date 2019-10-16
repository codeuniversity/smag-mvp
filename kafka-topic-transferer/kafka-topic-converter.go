package kafka_topic_transferer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/segmentio/kafka-go"
)

// Transferer represents the Transferer containing all clients it uses
type Transferer struct {
	kafkaTopicIn  *kafka.Reader
	kafkaTopicOut *kafka.Writer

	*service.Executor
}

// New returns an initilized Transferer
func New(topicIn *kafka.Reader, topicOut *kafka.Writer) *Transferer {
	c := &Transferer{}
	c.kafkaTopicIn = topicIn
	c.kafkaTopicOut = topicOut

	c.Executor = service.New()

	return c
}

//Run the Transferer
func (c *Transferer) Run() {
	fmt.Println("Start Transferer")
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
