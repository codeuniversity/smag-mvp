package main

import (
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/elastic"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"strconv"

	elasticsearch_inserter "github.com/codeuniversity/smag-mvp/elastic/indexer"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")

	esHosts := utils.GetMultipliesStringsFromEnvDefault("ELASTIC_SEARCH_ADDRESS", []string{"localhost:9201"})

	elasticInserter := elasticsearch_inserter.New(esHosts, elastic.PostsIndex, elastic.PostsIndexMapping, kafkaAddress, changesTopic, groupID, indexPost)

	service.CloseOnSignal(elasticInserter)
	waitUntilClosed := elasticInserter.Start()

	waitUntilClosed()
}

func indexPost(client *elasticsearch.Client, m *changestream.ChangeMessage) error {
	currentPost := &post{}
	err := json.Unmarshal(m.Payload.After, currentPost)

	if err != nil {
		return err
	}

	switch m.Payload.Op {
	case "r", "c":
		return upsertPost(currentPost, client)
	case "u":
		previousPost := &post{}
		err := json.Unmarshal(m.Payload.Before, previousPost)

		if err != nil {
			return err
		}

		if previousPost.Caption != currentPost.Caption {
			return upsertPost(currentPost, client)
		}
	}

	return nil
}

type post struct {
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	Caption string `json:"caption"`
}

func upsertPost(post *post, client *elasticsearch.Client) error {

	upsertBody := createUpsertBody(post)
	response, err := client.Update(elastic.PostsIndex, strconv.Itoa(post.ID), esutil.NewJSONReader(upsertBody))

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("upsertPost Upsert Document Failed StatusCode=%s Body=%s", response.Status(), response.String())
	}
	return nil
}

func createUpsertBody(post *post) map[string]interface{} {
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

	return commentUpsert
}
