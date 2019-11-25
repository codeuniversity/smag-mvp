package main

import (
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	elasticsearch_inserter "github.com/codeuniversity/smag-mvp/kafka/changestream/elastic-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/elastic/go-elasticsearch/v7"
	"strconv"
	"strings"
)

const instaPostUpdate = `
{
    "script" : {
        "source": "ctx._source.caption = params.caption",
        "lang": "painless",
        "params" : {
            "caption" : "%s"
        }
    }
}
`
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
		"userId": "%s"
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
        "userId" : {
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

	elasticInserter := elasticsearch_inserter.New(esHosts, esIndex, instaPostMapping, kafkaAddress, changesTopic, groupID, insertPost)

	service.CloseOnSignal(elasticInserter)
	waitUntilClosed := elasticInserter.Start()

	waitUntilClosed()
}

func insertPost(m *changestream.ChangeMessage, client *elasticsearch.Client) error {
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

		if previousPost.Caption != previousPost.Caption {
			return upsertPost(currentPost, client)
		}
	}

	return nil
}

type post struct {
	ID      int    `json:"id"`
	UserId  string `json:"user_id"`
	Caption string `json:"caption"`
}

func upsertPost(post *post, client *elasticsearch.Client) error {
	instaComment := fmt.Sprintf(instaPostUpsert, post.Caption, post.UserId, post.Caption)
	response, err := client.Update("insta_comments", strconv.Itoa(post.ID), strings.NewReader(instaComment))

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("FindPostId Update Document Failed StatusCode: %d", response.StatusCode)
	}
	return nil
}
