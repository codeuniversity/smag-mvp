package transferer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"

	"github.com/segmentio/kafka-go"
)

// Transferer represents the Transferer containing all clients it uses
type Transferer struct {
	fromTopic *kafka.Reader
	toTopic   *kafka.Writer

	*worker.Worker
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

	b := worker.Builder{}.WithName("kafka-topic-transferer").
		WithWorkStep(t.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("fromTopic", t.fromTopic.Close).
		AddShutdownHook("toTopic", t.toTopic.Close)

	t.Worker = b.MustBuild()

	return t
}

//Run the Transferer
func (t *Transferer) runStep() error {
	fmt.Println("Start Transferer")
	m, err := t.fromTopic.FetchMessage(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	changeMessage := &models.ChangeMessage{}
	err = json.Unmarshal(m.Value, changeMessage)
	if err != nil {
		fmt.Println(err)
		return err
	}
	userName := changeMessage.Payload.After.UserName
	t.fromTopic.CommitMessages(context.Background(), m)

	// checks if user was created or updated
	if changeMessage.Payload.Op == "c" {
		t.toTopic.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(userName),
		})

		log.Printf("%s tranfered \n", userName)

	} else {
		log.Printf("%s updated \n", userName)
	}
	return nil

}
