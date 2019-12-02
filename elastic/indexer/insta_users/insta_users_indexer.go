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
	bulkSize := utils.GetNumberFromEnvWithDefault("BULK_SIZE", 10)
	esHosts := utils.GetMultipleStringsFromEnvWithDefault("ES_HOSTS", []string{"http://localhost:9201"})

	i := indexer.New(esHosts, elastic.UsersIndex, elastic.UsersIndexMapping, kafkaAddress, changesTopic, groupID, handleChangeMessage, bulkSize)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

// handleChangeMessage filters relevant events and upserts them
func handleChangeMessage(m *changestream.ChangeMessage) (*indexer.ElasticIndexer, error) {
	user := &models.InstaUser{}
	if err := json.Unmarshal(m.Payload.After, user); err != nil {
		return &indexer.ElasticIndexer{}, err
	}

	switch m.Payload.Op {
	case "c", "r", "u":
		return createBulkUpsertOperation(user)
	}
	return &indexer.ElasticIndexer{}, nil
}

func createBulkUpsertOperation(user *models.InstaUser) (*indexer.ElasticIndexer, error) {
	var bulkOperation = map[string]interface{}{
		"update": map[string]interface{}{
			"_id":    user.ID,
			"_index": elastic.UsersIndex,
		},
	}

	bulkOperationJson, err := json.Marshal(bulkOperation)

	if err != nil {
		return &indexer.ElasticIndexer{}, err
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
		return &indexer.ElasticIndexer{}, err
	}

	usersUpsertJson = append(usersUpsertJson, "\n"...)

	bulkUpsertBody := string(bulkOperationJson) + string(usersUpsertJson)

	return &indexer.ElasticIndexer{DocumentId: strconv.Itoa(user.ID), BulkOperation: bulkUpsertBody}, err
}
