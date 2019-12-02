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
	bulkChunkSize := utils.GetNumberFromEnvWithDefault("BULK_CHUNK_SIZE", 10)
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")

	esHosts := utils.GetMultipleStringsFromEnvWithDefault("ES_HOSTS", []string{"localhost:9201"})

	i := indexer.New(esHosts, elastic.CommentsIndex, elastic.CommentsIndexMapping, kafkaAddress, changesTopic, groupID, indexComment, bulkChunkSize)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

func indexComment(m *changestream.ChangeMessage) (*indexer.ElasticIndexer, error) {
	comment := &models.InstaComment{}
	err := json.Unmarshal(m.Payload.After, comment)

	if err != nil {
		return &indexer.ElasticIndexer{}, err
	}

	switch m.Payload.Op {
	case "r", "c":
		return createBulkUpsertOperation(comment)
	}

	return &indexer.ElasticIndexer{}, nil
}

func createBulkUpsertOperation(comment *models.InstaComment) (*indexer.ElasticIndexer, error) {
	var bulkOperation = map[string]interface{}{
		"update": map[string]interface{}{
			"_id":    comment.ID,
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
			"source": "ctx._source.comment = params.comment",
			"lang":   "painless",
			"params": map[string]interface{}{
				"comment": comment.Comment,
			},
		},
		"upsert": map[string]interface{}{
			"post_id": comment.PostID,
			"comment": comment.Comment,
		},
	}

	commentUpsertJson, err := json.Marshal(commentUpsert)

	if err != nil {
		return &indexer.ElasticIndexer{}, err
	}

	commentUpsertJson = append(commentUpsertJson, "\n"...)

	bulkUpsertBody := string(bulkOperationJson) + string(commentUpsertJson)

	return &indexer.ElasticIndexer{DocumentId: strconv.Itoa(comment.ID), BulkOperation: bulkUpsertBody}, err
}
