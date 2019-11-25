package elasticsearch_inserter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/segmentio/kafka-go"
)

const (
	elasticProtocol = "http://%s"
)

// Inserter is the type definition of the esInserter
type Inserter struct {
	*worker.Worker

	esClient *elasticsearch.Client
	kReader  *kafka.Reader

	insertFunc InserterFunc
}

// InserterFunc is the type for the functions which will insert data into elasticsearch
type InserterFunc func(*changestream.ChangeMessage, *elasticsearch.Client) error

// New return a set up esInserter
func New(esHosts []string, esIndex, esMapping, kafkaAddress, changesTopic, kafkaGroupID string, inserterFunc InserterFunc) *Inserter {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)

	i := &Inserter{}
	i.kReader = kf.NewReader(readerConfig)
	i.insertFunc = inserterFunc

	i.esClient = i.initializeElasticSearch(esHosts)

	i.Worker = worker.Builder{}.WithName("elasticsearch-inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10 * time.Second).
		MustBuild()

	utils.PanicIfNotNil(i.createIndex(esIndex, esMapping))
	return i
}

func (i *Inserter) runStep() error {
	m, err := i.kReader.FetchMessage(context.Background())

	if err != nil {
		return err
	}

	changeMessage := &changestream.ChangeMessage{}

	err = json.Unmarshal(m.Value, changeMessage)

	if err != nil {
		return err
	}

	err = i.insertFunc(changeMessage, i.esClient)

	if err != nil {
		return err
	}

	log.Println("Inserted")
	return i.kReader.CommitMessages(context.Background(), m)
}

func (i *Inserter) initializeElasticSearch(esHosts []string) *elasticsearch.Client {

	var hosts []string
	for _, address := range esHosts {
		url := fmt.Sprintf(elasticProtocol, address)
		hosts = append(hosts, url)
	}

	cfg := elasticsearch.Config{
		Addresses: hosts,
	}
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		panic(err)
	}
	return client
}

func (i *Inserter) createIndex(esIndex, esMapping string) error {
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
		log.Println(response)
	} else if response.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("error finding index: %d %s", response.StatusCode, string(body))
	}
	return nil
}
