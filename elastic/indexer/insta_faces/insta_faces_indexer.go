package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"

	"github.com/codeuniversity/smag-mvp/elastic"
	"github.com/codeuniversity/smag-mvp/elastic/indexer"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/models"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.face_data")

	esHosts := utils.GetMultipliesStringsFromEnvDefault("ES_HOSTS", []string{"http://localhost:9200"})

	i := indexer.New(esHosts, elastic.FacesIndex, elastic.FacesIndexMapping, kafkaAddress, changesTopic, groupID, indexFace)

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

	doc, err := FaceDocFromFaceData(face)
	if err != nil {
		return err
	}

	docReader := esutil.NewJSONReader(doc)
	response, err := client.Index(elastic.FacesIndex,
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
