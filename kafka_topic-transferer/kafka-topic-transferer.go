package transferer

import (
	"context"
	"encoding/json"
	"fmt"

	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"

	"github.com/segmentio/kafka-go"
)

// Transferer represents the Transferer containing all clients it uses
type Transferer struct {
	fromTopic *kafka.Reader
	toTopic   *kafka.Writer

	*service.Executor
}

// New returns an initilized Transferer
func New(fromTopic string, toTopic string) *Transferer {
	c := &Transferer{}

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	readerConfig := kf.NewReaderConfig(kafkaAddress, groupID, fromTopic)
	writerConfig := kf.NewWriterConfig(kafkaAddress, toTopic, true)

	c.fromTopic = kf.NewReader(readerConfig)
	c.toTopic = kf.NewWriter(writerConfig)

	c.Executor = service.New()

	return c
}

//Run the Transferer
func (c *Transferer) Run() {
	fmt.Println("Start Transferer")
	for c.IsRunning() {
		m, err := c.fromTopic.FetchMessage(context.Background())
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
		c.fromTopic.CommitMessages(context.Background(), m)

		c.toTopic.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(userName),
		})

		fmt.Printf("%s tranfered \n", userName)

	}
}
