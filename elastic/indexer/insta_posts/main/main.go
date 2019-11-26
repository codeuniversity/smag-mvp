package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"

	elasticsearch_inserter "github.com/codeuniversity/smag-mvp/elastic/indexer"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

const instaPostUpsert = `
	{
    "script" : {
        "source": "ctx._source.caption = params.caption",
        "lang": "painless",
        "params" : {
            "caption" : %s
        }
    },
    "upsert" : {
		"user_id": "%s"
		"caption": "%s"
    }
}
`

const instaPostMapping = `
	{
    "mappings" : {
      "properties" : {
        "caption" : {
          "type" : "text"
        },
        "user_id" : {
          "type" : "keyword"
        }
      }
    }
}
`

const esIndex = "insta_posts"

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")

	esHosts := utils.GetMultipliesStringsFromEnvDefault("ELASTIC_SEARCH_ADDRESS", []string{"localhost:9201"})

	elasticInserter := elasticsearch_inserter.New(esHosts, esIndex, instaPostMapping, kafkaAddress, changesTopic, groupID, indexPost)

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
	instaComment := fmt.Sprintf(instaPostUpsert, post.Caption, post.UserID, post.Caption)
	response, err := client.Update(esIndex, strconv.Itoa(post.ID), strings.NewReader(instaComment))

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("upsertPost Upsert Document Failed StatusCode=%s Body=%s", response.Status(), response.String())
	}
	return nil
}
