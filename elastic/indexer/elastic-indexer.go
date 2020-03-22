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

	esClient                *elasticsearch.Client
	kReader                 *kafka.Reader
	esIndex                 string
	bulkChunkSize           int
	bulkFetchTimeoutSeconds int
	indexFunc               IndexFunc
}

// IndexFunc is the type for the functions which will insert data into elasticsearch
type IndexFunc func(*changestream.ChangeMessage) (*BulkIndexDoc, error)

// New returns an initialised Indexer
func New(esHosts []string, esIndex, esMapping, kafkaAddress, changesTopic, kafkaGroupID string, indexFunc IndexFunc, bulkChunkSize int, bulkFetchTimeout int) *Indexer {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)

	i := &Indexer{}
	i.kReader = kf.NewReader(readerConfig)
	i.indexFunc = indexFunc
	i.esIndex = esIndex
	i.esClient = elastic.InitializeElasticSearch(esHosts)
	i.bulkChunkSize = bulkChunkSize
	i.bulkFetchTimeoutSeconds = bulkFetchTimeout

	i.Worker = worker.Builder{}.WithName(fmt.Sprintf("indexer[%s->es/%s]", changesTopic, esIndex)).
		WithWorkStep(i.runStep).
		WithStopTimeout(10 * time.Second).
		MustBuild()

	utils.PanicIfNotNil(i.createIndex(esIndex, esMapping))
	return i
}

func (i *Indexer) runStep() error {
	messages, err := i.readMessageBlock(i.bulkChunkSize)
	log.Println("Messages Bulk: ", len(messages))
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	var bulkBody string
	bulkDocumentIdKafkaMessages := make(map[string]kafka.Message)
	for _, message := range messages {

		changeMessage := &changestream.ChangeMessage{}
		if err := json.Unmarshal(message.Value, changeMessage); err != nil {
			return err
		}
		bulkOperation, err := i.indexFunc(changeMessage)

		if err != nil {
			return err
		}

		bulkDocumentIdKafkaMessages[bulkOperation.DocumentId] = message
		if bulkOperation == nil {
			continue
		}
		bulkBody += bulkOperation.BulkOperation
	}

	log.Println("RequestBody: ", bulkBody)
	if bulkBody == "" {
		err := i.kReader.CommitMessages(context.Background(), messages...)
		if err != nil {
			return err
		}
		return nil
	}
	bulkResponse, err := i.esClient.Bulk(strings.NewReader(bulkBody), i.esClient.Bulk.WithIndex(i.esIndex))
	if err != nil {
		return err
	}
	log.Println("Result Messages Bulk: ", bulkResponse.Status())

	body, err := ioutil.ReadAll(bulkResponse.Body)
	if err != nil {
		return err
	}

	var result bulkResult
	err = json.Unmarshal(body, &result)

	if err != nil {
		return err
	}

	log.Println("BulkResultItem: ", len(result.Items))
	log.Println("BulResult: ", string(body))
	err = i.checkAllResultMessagesAreValid(&result)
	if err != nil {
		return err
	}
	for _, message := range bulkDocumentIdKafkaMessages {
		err := i.kReader.CommitMessages(context.Background(), message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Indexer) checkAllResultMessagesAreValid(result *bulkResult) error {
	if result == nil {
		return fmt.Errorf("BulkResult is nil")
	}

	for _, bulkResultOperation := range result.Items {
		if bulkResultOperation.Index != nil {

			err := errorForHttpStatus(bulkResultOperation.Index.Status)
			if err != nil {
				return err
			}
		} else if bulkResultOperation.Update != nil {

			err := errorForHttpStatus(bulkResultOperation.Update.Status)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func errorForHttpStatus(httpStatus int) error {
	if httpStatus != 200 && httpStatus != 201 {
		return fmt.Errorf("creating/updateing index failed Httpstatus: %d \n", httpStatus)
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

func (i *Indexer) readMessageBlock(maxChunkSize int) (messages []kafka.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(i.bulkFetchTimeoutSeconds))
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
		Index *struct {
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
		} `json:"index,omitempty"`
		Update *struct {
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
		} `json:"update,omitempty"`
	} `json:"items"`
}
