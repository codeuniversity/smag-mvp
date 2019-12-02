package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
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

	esClient  *elasticsearch.Client
	kReader   *kafka.Reader
	esIndex   string
	indexFunc IndexFunc
}

const bulkLimit = 200

// IndexFunc is the type for the functions which will insert data into elasticsearch
type IndexFunc func(*changestream.ChangeMessage) (*ElasticIndexer, error)

// New returns an initialised Indexer
func New(esHosts []string, esIndex, esMapping, kafkaAddress, changesTopic, kafkaGroupID string, indexFunc IndexFunc) *Indexer {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)

	i := &Indexer{}
	i.kReader = kf.NewReader(readerConfig)
	i.indexFunc = indexFunc
	i.esIndex = esIndex
	i.esClient = elastic.InitializeElasticSearch(esHosts)

	i.Worker = worker.Builder{}.WithName(fmt.Sprintf("indexer[%s->es/%s]", changesTopic, esIndex)).
		WithWorkStep(i.runStep).
		WithStopTimeout(10 * time.Second).
		MustBuild()

	utils.PanicIfNotNil(i.createIndex(esIndex, esMapping))
	return i
}

func (i *Indexer) runStep() error {
	messages, err := i.readMessageBlock(1*time.Second, bulkLimit)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	var bulkBody string
	bulkDocumentIdMessages := make(map[string]kafka.Message)
	for _, message := range messages {

		changeMessage := &changestream.ChangeMessage{}
		if err := json.Unmarshal(message.Value, changeMessage); err != nil {
			return err
		}
		bulkOperation, err := i.indexFunc(changeMessage)

		if err != nil {
			return err
		}

		bulkDocumentIdMessages[bulkOperation.DocumentId] = message
		bulkBody += bulkOperation.BulkOperation
	}

	bulkResponse, err := i.esClient.Bulk(strings.NewReader(bulkBody), i.esClient.Bulk.WithIndex(i.esIndex))

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(bulkResponse.Body)

	if err != nil {
		return err
	}

	var bulkResult bulkResult
	err = json.Unmarshal(body, &bulkResult)

	if err != nil {
		return err
	}

	for _, bulkResultOperation := range bulkResult.Items {

		if bulkResultOperation.Index.Status == 200 || bulkResultOperation.Index.Status == 201 {
			err := i.kReader.CommitMessages(context.Background(), bulkDocumentIdMessages[bulkResultOperation.Index.ID])

			if err != nil {
				return err
			}
		}
	}
	return nil
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

func (i *Indexer) readMessageBlock(timeout time.Duration, maxChunkSize int) (messages []kafka.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for k := 0; k < maxChunkSize; k++ {
		m, err := i.kReader.FetchMessage(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				return messages, nil
			}

			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}

type bulkResult struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			Index   string `json:"_index"`
			Type    string `json:"_type"`
			ID      string `json:"_id"`
			Version int    `json:"_version"`
			Result  string `json:"result"`
			Shards  struct {
				Total      int `json:"total"`
				Successful int `json:"successful"`
				Failed     int `json:"failed"`
			} `json:"_shards"`
			SeqNo       int `json:"_seq_no"`
			PrimaryTerm int `json:"_primary_term"`
			Status      int `json:"status"`
		} `json:"index"`
	} `json:"items"`
}
