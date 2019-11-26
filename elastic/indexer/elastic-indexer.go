package elasticsearch_inserter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/codeuniversity/smag-mvp/elastic"

	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/segmentio/kafka-go"
)

// Indexer is the type definition of the esInserter
type Indexer struct {
	*worker.Worker

	esClient *elasticsearch.Client
	kReader  *kafka.Reader

	indexFunc IndexFunc
}

// IndexFunc is the type for the functions which will insert data into elasticsearch
type IndexFunc func(*elasticsearch.Client, *changestream.ChangeMessage) error

// New returns an initialised Indexer
func New(esHosts []string, esIndex, esMapping, kafkaAddress, changesTopic, kafkaGroupID string, indexFunc IndexFunc) *Indexer {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)

	i := &Indexer{}
	i.kReader = kf.NewReader(readerConfig)
	i.indexFunc = indexFunc

	i.esClient = elastic.InitializeElasticSearch(esHosts)

	i.Worker = worker.Builder{}.WithName(fmt.Sprintf("indexer[%s->es/%s]", changesTopic, esIndex)).
		WithWorkStep(i.runStep).
		WithStopTimeout(10 * time.Second).
		MustBuild()

	utils.PanicIfNotNil(i.createIndex(esIndex, esMapping))
	return i
}

func (i *Indexer) runStep() error {
	m, err := i.kReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	changeMessage := &changestream.ChangeMessage{}
	if err := json.Unmarshal(m.Value, changeMessage); err != nil {
		return err
	}

	if err := i.indexFunc(i.esClient, changeMessage); err != nil {
		return err
	}
	return i.kReader.CommitMessages(context.Background(), m)
}

func (i *Indexer) createIndex(esIndex, esMapping string) error {
	response, err := i.esClient.Indices.Exists(
		[]string{esIndex},
		i.esClient.Indices.Exists.WithHuman(),
		i.esClient.Indices.Exists.WithPretty(),
	)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode == 404 {
		bodyReader := bytes.NewReader([]byte(esMapping))
		response, err := i.esClient.Indices.Create(
			esIndex,
			i.esClient.Indices.Create.WithHuman(),
			i.esClient.Indices.Create.WithPretty(),
			i.esClient.Indices.Create.WithBody(bodyReader),
		)

		if err != nil {
			return err
		}
		log.Println(response.String())
	} else if response.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("error finding index: %d %s", response.StatusCode, string(body))
	}
	return nil
}
