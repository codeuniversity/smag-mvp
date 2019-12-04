package main

import (
	"encoding/json"
	"strconv"

	"github.com/codeuniversity/smag-mvp/elastic"
	"github.com/codeuniversity/smag-mvp/elastic/indexer"
	"github.com/codeuniversity/smag-mvp/elastic/models"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.GetStringFromEnvWithDefault("KAFKA_GROUPID", "insta_usersearch-inserter")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.users")
	bulkChunkSize := utils.GetNumberFromEnvWithDefault("BULK_CHUNK_SIZE", 10)
	bulkFetchTimeoutSeconds := utils.GetNumberFromEnvWithDefault("BULK_FETCH_TIMEOUT_SECONDS", 5)
	esHosts := utils.GetMultipleStringsFromEnvWithDefault("ES_HOSTS", []string{"http://localhost:9201"})

	i := indexer.New(esHosts, elastic.UsersIndex, elastic.UsersIndexMapping, kafkaAddress, changesTopic, groupID, handleChangeMessage, bulkChunkSize, bulkFetchTimeoutSeconds)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

// handleChangeMessage filters relevant events and upserts them
func handleChangeMessage(m *changestream.ChangeMessage) (*indexer.BulkIndexDoc, error) {
	user := &models.InstaUser{}
	if err := json.Unmarshal(m.Payload.After, user); err != nil {
		return nil, err
	}

	switch m.Payload.Op {
	case "c", "r", "u":
		return createBulkUpsertOperation(user)
	}
	return nil, nil
}

func createBulkUpsertOperation(user *models.InstaUser) (*indexer.BulkIndexDoc, error) {
	var bulkOperation = map[string]interface{}{
		"update": map[string]interface{}{
			"_id":    user.ID,
			"_index": elastic.UsersIndex,
		},
	}

	bulkOperationJson, err := json.Marshal(bulkOperation)

	if err != nil {
		return nil, err
	}

	bulkOperationJson = append(bulkOperationJson, "\n"...)

	var usersUpsert = map[string]interface{}{
		"script": map[string]interface{}{
			"source": "ctx._source.user_name = params.user_name; ctx._source.real_name = params.real_name; ctx._source.bio = params.bio",
			"lang":   "painless",
			"params": map[string]interface{}{
				"user_name": user.Username,
				"real_name": user.Realname,
				"bio":       user.Bio,
			},
		},
		"upsert": map[string]interface{}{
			"user_name": user.Username,
			"real_name": user.Realname,
			"bio":       user.Bio,
		},
	}

	usersUpsertJson, err := json.Marshal(usersUpsert)

	if err != nil {
		return nil, err
	}

	usersUpsertJson = append(usersUpsertJson, "\n"...)

	bulkUpsertBody := string(bulkOperationJson) + string(usersUpsertJson)

	return &indexer.BulkIndexDoc{DocumentId: strconv.Itoa(user.ID), BulkOperation: bulkUpsertBody}, err
}
