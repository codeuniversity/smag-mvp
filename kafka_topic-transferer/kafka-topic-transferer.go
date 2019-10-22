package transferer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	t := &Transferer{}

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	readerConfig := kf.NewReaderConfig(kafkaAddress, groupID, fromTopic)
	writerConfig := kf.NewWriterConfig(kafkaAddress, toTopic, true)

	t.fromTopic = kf.NewReader(readerConfig)
	t.toTopic = kf.NewWriter(writerConfig)

	t.Executor = service.New()

	return t
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

		// checks if user was created or updated
		if changeStream.Payload.Op == "c" {
			c.toTopic.WriteMessages(context.Background(), kafka.Message{
				Value: []byte(userName),
			})

			fmt.Printf("%s tranfered \n", userName)

		} else {
			fmt.Printf("%s updated \n", userName)
		}

	}
}

//Close the transferer
func (t *Transferer) Close() {
	t.Stop()
	t.WaitUntilStopped(time.Second * 3)

	t.fromTopic.Close()
	t.toTopic.Close()

	t.MarkAsClosed()
}
