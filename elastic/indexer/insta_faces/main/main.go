package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"

	elasticsearch_indexer "github.com/codeuniversity/smag-mvp/elastic/indexer"
	indexer "github.com/codeuniversity/smag-mvp/elastic/indexer/insta_faces"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

const indexCreateBody = `
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
const esIndex = "faces"

func main() {
	esHosts := utils.GetMultipliesStringsFromEnvDefault("ES_HOSTS", []string{"http://localhost:9200"})
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.face_data")

	i := elasticsearch_indexer.New(esHosts, esIndex, indexCreateBody, kafkaAddress, changesTopic, groupID, indexFace)

	service.CloseOnSignal(i)
	waitUntilDone := i.Start()

	waitUntilDone()
}

func indexFace(client *elasticsearch.Client, m *changestream.ChangeMessage) error {

	switch m.Payload.Op {
	case "r", "u", "c":
		break
	default:
		return nil
	}
	face := &models.FaceData{}
	err := json.Unmarshal(m.Payload.After, face)
	if err != nil {
		return err
	}

	doc, err := indexer.FaceDocFromFaceData(face)
	if err != nil {
		return err
	}

	docReader := esutil.NewJSONReader(doc)
	response, err := client.Index(esIndex,
		docReader,
		client.Index.WithHuman(),
		client.Index.WithPretty(),
	)
	if err != nil {
		return err
	}
	log.Println(response.StatusCode, response.String())

	if response.StatusCode != 201 {
		return fmt.Errorf("IndexFace Create Document Failed StatusCode=%s Body=%s", response.Status(), response.String())
	}

	return nil
}
