package filter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"

	"github.com/segmentio/kafka-go"
)

// Filter represents the Filter containing all clients it uses
type Filter struct {
	*worker.Worker

	fromTopic *kafka.Reader
	toTopic   *kafka.Writer
}

// New returns an initilized Filter
func New(fromTopic string, toTopic string) *Filter {
	t := &Filter{}

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	readerConfig := kf.NewReaderConfig(kafkaAddress, groupID, fromTopic)
	writerConfig := kf.NewWriterConfig(kafkaAddress, toTopic, true)

	t.fromTopic = kf.NewReader(readerConfig)
	t.toTopic = kf.NewWriter(writerConfig)

	b := worker.Builder{}.WithName("kafka_topic_transferer").
		WithWorkStep(t.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("fromTopic", t.fromTopic.Close).
		AddShutdownHook("toTopic", t.toTopic.Close)

	t.Worker = b.MustBuild()

	return t
}

func (t *Filter) runStep() error {
	m, err := t.fromTopic.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	changeMessage := &changeMessage{}
	err = json.Unmarshal(m.Value, changeMessage)
	if err != nil {
		fmt.Println(err)
		return err
	}
	userName := changeMessage.Payload.After.UserName
	t.fromTopic.CommitMessages(context.Background(), m)

	// checks if user was created
	if changeMessage.Payload.Op == "c" {
		t.toTopic.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(userName),
		})

		log.Println(userName, "transfered")

	} else {
		log.Println(userName, "updated")
	}
	return nil

}
