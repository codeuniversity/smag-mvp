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
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	bulkSize := utils.GetNumberFromEnvWithDefault("BULK_SIZE", 10)
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")

	esHosts := utils.GetMultipleStringsFromEnvWithDefault("ES_HOSTS", []string{"localhost:9201"})

	i := indexer.New(esHosts, elastic.PostsIndex, elastic.PostsIndexMapping, kafkaAddress, changesTopic, groupID, indexPost, bulkSize)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

func indexPost(m *changestream.ChangeMessage) (*indexer.ElasticIndexer, error) {
	currentPost := &models.InstaPost{}
	err := json.Unmarshal(m.Payload.After, currentPost)

	if err != nil {
		return &indexer.ElasticIndexer{}, err
	}

	switch m.Payload.Op {
	case "r", "c":
		return createBulkUpsertOperation(currentPost)
	case "u":
		previousPost := &models.InstaPost{}
		err := json.Unmarshal(m.Payload.Before, previousPost)

		if err != nil {
			return &indexer.ElasticIndexer{}, err
		}

		if previousPost.Caption != currentPost.Caption {
			return createBulkUpsertOperation(currentPost)
		}
	}

	return &indexer.ElasticIndexer{}, nil
}

func createBulkUpsertOperation(post *models.InstaPost) (*indexer.ElasticIndexer, error) {
	var bulkOperation = map[string]interface{}{
		"update": map[string]interface{}{
			"_id":    post.ID,
			"_index": elastic.PostsIndex,
		},
	}

	bulkOperationJson, err := json.Marshal(bulkOperation)

	if err != nil {
		return &indexer.ElasticIndexer{}, err
	}

	bulkOperationJson = append(bulkOperationJson, "\n"...)

	var commentUpsert = map[string]interface{}{
		"script": map[string]interface{}{
			"source": "ctx._source.caption = params.caption",
			"lang":   "painless",
			"params": map[string]interface{}{
				"caption": post.Caption,
			},
		},
		"upsert": map[string]interface{}{
			"user_id": post.UserID,
			"caption": post.Caption,
		},
	}

	commentUpsertJson, err := json.Marshal(commentUpsert)

	if err != nil {
		return &indexer.ElasticIndexer{}, err
	}

	commentUpsertJson = append(commentUpsertJson, "\n"...)

	bulkUpsertBody := string(bulkOperationJson) + string(commentUpsertJson)

	return &indexer.ElasticIndexer{DocumentId: strconv.Itoa(post.ID), BulkOperation: bulkUpsertBody}, err
}
