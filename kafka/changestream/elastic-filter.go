package changestream

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/segmentio/kafka-go"
)

// KafkaToElasticFilter represents the scraper containing all clients it uses
type KafkaToElasticFilter struct {
	*worker.Worker

	changesReader *kafka.Reader

	// TODO: Add elasticWriter to KafkaToElasticFilter
	// elasticWriter connection
}

// New returns an initilized KafkaToElasticFilter
func New(kafkaAddress, kafkaGroupID, changesTopic, elasticTopic string) *KafkaToElasticFilter {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)

	f := &KafkaToElasticFilter{
		changesReader: kf.NewReader(readerConfig),
	}

	b := worker.Builder{}.
		WithName(fmt.Sprintf("KafkaToElasticFilter[%s->%s]", changesTopic, elasticTopic)).
		WithWorkStep(f.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("changesReader", f.changesReader.Close) //.
		// AddShutdownHook("elasticWriter", ...)
		// TODO: Add ShutdownHook for elasticWriter

	f.Worker = b.MustBuild()

	return f
}

func (f *KafkaToElasticFilter) runStep() error {
	m, err := f.changesReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	changeMessage := &ChangeMessage{}
	err = json.Unmarshal(m.Value, changeMessage)
	if err != nil {
		return err
	}

	elasticMessages, err := filterChange(changeMessage)
	if err != nil {
		return err
	}

	log.Println(elasticMessages)

	// TODO: Write elasticMessages to elasticsearch
	// if len(elasticMessages) > 0 {
	// 	err = elasticWriter.Write(...)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return f.changesReader.CommitMessages(context.Background(), m)
}

// FILTER
//

type user struct {
	ID       int    `json:"id"`
	Username string `json:"user_name"`
}

// TODO: Adapt return value to Elsatic search
func filterChange(m *ChangeMessage) ([]kafka.Message, error) {
	// use create (c) or update (u) events
	if (m.Payload.Op != "c") || (m.Payload.Op != "u") {
		return nil, nil
	}

	// unmarshal into user to get only relevant data
	u := &user{}
	err := json.Unmarshal(m.Payload.After, u)
	if err != nil {
		return nil, err
	}

	// marshal the relevant data again
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return []kafka.Message{{Value: b}}, nil
}
