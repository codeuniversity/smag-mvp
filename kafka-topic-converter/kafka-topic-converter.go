package kafka-topic-converter

import (
	"fmt"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/models"
)

type Converter struct {
	kafkaTopicIn *kafka.Reader
	kafkaTopicOut *kafka.Writer

	*service.Executor
}


func New(topicIn *kafa.Reader, topicOut *kafka.Writer) *Converter {
	c := &Converter{}
	c.kafkaTopicIn = topicIn
	c.kafkaTopicOut = topicOut

	c.Executor = service.New()

	return i
}


func (c *Converter) Run() {

	for c.IsRunning() {
		m, err := c.kafkaTopicIn.FetchMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}

		debeziumModel := &models.DebeziumTopic{}
		err = json.Unmarshal(m,debeziumModel)
		fmt.Println(debeziumModel)

	}
}
