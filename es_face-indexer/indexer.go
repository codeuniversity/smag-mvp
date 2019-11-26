package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/codeuniversity/smag-mvp/es"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/segmentio/kafka-go"
)

// Indexer indexes faces in elasticsearch with their encoding
// creates the index and mapping if the it doesn't on startup
type Indexer struct {
	changeQReader *kafka.Reader
	*worker.Worker
	esClient *elasticsearch.Client
}

// New returns an initialized Indexer and creates the Index in elasticsearch if it doesn't exist yet
func New(changeQReader *kafka.Reader, esHosts []string) *Indexer {
	cfg := elasticsearch.Config{Addresses: esHosts}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating elasticsearch client: %s", err)
	}
	i := &Indexer{
		changeQReader: changeQReader,
		esClient:      esClient,
	}

	i.Worker = worker.Builder{}.WithName("face-indexer").
		WithWorkStep(i.step).
		AddShutdownHook("changeQReader", changeQReader.Close).
		MustBuild()

	utils.PanicIfNotNil(i.createIndex())

	return i
}

var indexCreateBody string = `
{
	"mappings" : {
			"properties" : {
				"encoding_vector": {
					"type": "binary",
					"doc_values": true
				},
				"post_id": {
					"type": "integer"
				},
				"x": {
					"type": "integer"
				},
				"y": {
					"type": "integer"
				},
				"width": {
					"type": "integer"
				},
				"height":{
					"type": "integer"
				}
			}
	}
}
`

func (i *Indexer) createIndex() error {
	response, err := i.esClient.Indices.Exists(
		[]string{es.FaceIndexName},
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
		bodyReader := bytes.NewReader([]byte(indexCreateBody))
		response, err := i.esClient.Indices.Create(
			es.FaceIndexName,
			i.esClient.Indices.Create.WithHuman(),
			i.esClient.Indices.Create.WithPretty(),
			i.esClient.Indices.Create.WithBody(bodyReader),
		)

		if err != nil {
			return err
		}
		fmt.Println(response)
	} else if response.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("error finding index: %d %s", response.StatusCode, string(body))
	}
	return nil
}

func (i *Indexer) step() error {
	m, err := i.changeQReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	change := &changestream.ChangeMessage{}
	err = json.Unmarshal(m.Value, change)
	if err != nil {
		return err
	}
	if change.Payload.Op == "d" {
		return i.changeQReader.CommitMessages(context.Background(), m)
	}
	face := &models.FaceData{}
	err = json.Unmarshal(change.Payload.After, face)
	if err != nil {
		return err
	}

	doc, err := FaceDocFromFaceData(face)
	if err != nil {
		return err
	}

	docReader := esutil.NewJSONReader(doc)
	response, err := i.esClient.Index(es.FaceIndexName,
		docReader,
		i.esClient.Index.WithHuman(),
		i.esClient.Index.WithPretty(),
	)
	if err != nil {
		return err
	}

	log.Println(response.StatusCode, response.String())

	return i.changeQReader.CommitMessages(context.Background(), m)
}
