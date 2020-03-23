package main

import (
	"encoding/json"
	"strconv"

	"github.com/codeuniversity/smag-mvp/elastic"
	"github.com/codeuniversity/smag-mvp/elastic/indexer"
	esModels "github.com/codeuniversity/smag-mvp/elastic/models"
	"github.com/codeuniversity/smag-mvp/insta/models"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	bulkChunkSize := utils.GetNumberFromEnvWithDefault("BULK_CHUNK_SIZE", 10)
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.face_data")
	bulkFetchTimeoutSeconds := utils.GetNumberFromEnvWithDefault("BULK_FETCH_TIMEOUT_SECONDS", 5)
	esHosts := utils.GetMultipleStringsFromEnvWithDefault("ES_HOSTS", []string{"http://localhost:9200"})

	i := indexer.New(esHosts, elastic.FacesIndex, elastic.FacesIndexMapping, kafkaAddress, changesTopic, groupID, indexFace, bulkChunkSize, bulkFetchTimeoutSeconds)

	service.CloseOnSignal(i)
	waitUntilDone := i.Start()

	waitUntilDone()
}

func indexFace(m *changestream.ChangeMessage) (*indexer.BulkIndexDoc, error) {

	switch m.Payload.Op {
	case "r", "u", "c":
		break
	default:
		return nil, nil
	}

	face := &models.FaceData{}
	err := json.Unmarshal(m.Payload.After, face)
	if err != nil {
		return nil, err
	}

	return createBulkIndexOperation(face)
}

func createBulkIndexOperation(face *models.FaceData) (*indexer.BulkIndexDoc, error) {
	bulkOperation := `{ "index": {}  }` + "\n"

	doc, err := esModels.FaceDocFromFaceData(face)
	if err != nil {
		return nil, err
	}

	docJson, err := json.Marshal(doc)

	if err != nil {
		return nil, err
	}

	docJson = append(docJson, "\n"...)

	bulkUpsertBody := bulkOperation + string(docJson)

	return &indexer.BulkIndexDoc{DocumentId: strconv.Itoa(int(face.ID)), BulkOperation: bulkUpsertBody}, err

}
