package changestream

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"

	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/segmentio/kafka-go"
)

// KafkaToElasticFilter represents the scraper containing all clients it uses
type KafkaToElasticFilter struct {
	*worker.Worker

	changesReader *kafka.Reader

	elasticWriter *elasticsearch.Client
}

// New returns an initilized KafkaToElasticFilter
func New(kafkaAddress, kafkaGroupID, changesTopic string, elasticAddresses []string) *KafkaToElasticFilter {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)

	f := &KafkaToElasticFilter{
		changesReader: kf.NewReader(readerConfig),
	}

	cfg := elasticsearch.Config{Addresses: elasticAddresses}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating elasticsearch client: %s", err)
	}
	f.elasticWriter = es

	b := worker.Builder{}.
		WithName(fmt.Sprintf("KafkaToElasticFilter[%s->%v]", changesTopic, elasticAddresses)).
		WithWorkStep(f.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("changesReader", f.changesReader.Close)
		// No shutdown hook for elasticWriter :(

	f.Worker = b.MustBuild()

	return f
}

func (f *KafkaToElasticFilter) runStep() error {
	m, err := f.changesReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	changeMessage := &ChangeMessage{}
	if err := json.Unmarshal(m.Value, changeMessage); err != nil {
		return err
	}

	elasticMessages, err := filterChange(changeMessage)
	if err != nil {
		return err
	}

	for _, userJSON := range elasticMessages {
		req := esapi.IndexRequest{
			Index: "user",
			// DocumentID: strconv.Itoa(i + 1), // not sure if it's needed to explicitly set the DocumentID
			Body:    strings.NewReader(userJSON),
			Refresh: "true",
		}
		res, err := req.Do(context.Background(), f.elasticWriter)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
	}

	return f.changesReader.CommitMessages(context.Background(), m)
}

// FILTER
//

type user struct {
	ID       int    `json:"id"`
	Username string `json:"user_name"`
}

// TODO: Adapt return value to Elsatic search
func filterChange(m *ChangeMessage) ([]string, error) {
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
	v, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	w := string(v)

	return []string{w}, nil
}
