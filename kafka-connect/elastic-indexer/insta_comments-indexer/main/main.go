package main

import (
	"encoding/json"
	"fmt"
	elasticsearch_inserter "github.com/codeuniversity/smag-mvp/kafka-connect/elastic-indexer"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/elastic/go-elasticsearch/v7"
	"strconv"
	"strings"
)

const instaCommentUpsert = `
	{
    "script" : {
        "source": "ctx._source.comment = params.comment",
        "lang": "painless",
        "params" : {
            "comment" : %s
        }
    },
    "upsert" : {
        "post_id" : "%s",
		"comment": "%s"
    }
}
`

const instaCommentsMapping = `
{
    "mappings" : {
      "properties" : {
        "comment" : {
          "type" : "text"
        },
        "post_id" : {
          "type" : "keyword"
        }
      }
    }
  }
`

const esIndex = "insta_comments"

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")

	esHosts := utils.GetMultipliesStringsFromEnvDefault("ELASTIC_SEARCH_ADDRESS", []string{"localhost:9201"})

	elasticInserter := elasticsearch_inserter.New(esHosts, esIndex, instaCommentsMapping, kafkaAddress, changesTopic, groupID, handleComment)

	service.CloseOnSignal(elasticInserter)
	waitUntilClosed := elasticInserter.Start()

	waitUntilClosed()
}

func handleComment(m *changestream.ChangeMessage, client *elasticsearch.Client) error {
	comment := &comment{}
	err := json.Unmarshal(m.Payload.After, comment)

	if err != nil {
		return err
	}

	switch m.Payload.Op {
	case "r", "c":
		return upsertComment(comment, client)
	}

	return nil
}

func upsertComment(comment *comment, client *elasticsearch.Client) error {
	instaComment := fmt.Sprintf(instaCommentUpsert, comment.Comment, comment.PostId, comment.Comment)
	response, err := client.Update(esIndex, strconv.Itoa(comment.ID), strings.NewReader(instaComment))

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("FindPostId Update Document Failed StatusCode: %d", response.StatusCode)
	}
	return nil
}

type comment struct {
	ID      int    `json:"id"`
	PostId  string `json:"post_id"`
	Comment string `json:"comment_text"`
}
