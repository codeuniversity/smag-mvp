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
        "postId" : "%s",
		"comment": "%s"
    }
}
`

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")

	elasticAddress := utils.MustGetStringFromEnv("ELASTIC_SEARCH_ADDRESS")

	elasticInserter := elasticsearch_inserter.New(elasticAddress, kafkaAddress, changesTopic, groupID, insertComment)

	service.CloseOnSignal(elasticInserter)
	waitUntilClosed := elasticInserter.Start()

	waitUntilClosed()
}

func insertComment(m *changestream.ChangeMessage, client *elasticsearch.Client) error {
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
	response, err := client.Update("insta_comments", strconv.Itoa(comment.ID), strings.NewReader(instaComment))

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
