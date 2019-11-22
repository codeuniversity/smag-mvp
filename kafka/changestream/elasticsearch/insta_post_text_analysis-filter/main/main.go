package main

import (
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	elasticsearch_inserter "github.com/codeuniversity/smag-mvp/kafka/changestream/elasticsearch"
	insta_post_text_filter "github.com/codeuniversity/smag-mvp/kafka/changestream/elasticsearch/insta_post_text_analysis-filter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/elastic/go-elasticsearch/v7"
	"io/ioutil"
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

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")

	elasticAddress := utils.MustGetStringFromEnv("ELASTIC_SEARCH_ADDRESS")

	inserter := elasticsearch_inserter.New(elasticAddress, kafkaAddress, changesTopic, groupID, insertPostText)

	service.CloseOnSignal(inserter)
	waitUntilClosed := inserter.Start()

	waitUntilClosed()
}

func insertPostText(m *changestream.ChangeMessage, client *elasticsearch.Client) error {
	post := &post{}
	err := json.Unmarshal(m.Payload.After, post)

	if err != nil {
		return err
	}

	switch m.Payload.Op {
	case "r", "c":
		return createPost(post, client)
	case "u":
		documentId, err := findPostId(strconv.Itoa(post.ID), client)

		if err != nil {
			return err
		}

		return updateDocumentId(documentId, post.Caption, client)
	}

	return nil
}

//these operator u, o to find documentId
func findPostId(postId string, client *elasticsearch.Client) (string, error) {

	postQuery := "postId: \"%s\""
	query := fmt.Sprintf(postQuery, postId)
	response, err := client.Search(client.Search.WithIndex("insta_post"), client.Search.WithQuery(query))

	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("FindPostId Search Query Failed StatusCode: %d", response.StatusCode)
	}

	searchResult := insta_post_text_filter.QueryResult{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &searchResult)

	if err != nil {
		panic(err)
	}

	if len(searchResult.Hits.Hits) == 1 {
		return "", fmt.Errorf("found Duplicate Documents: PostId: %s Duplicates: %d", postId, len(searchResult.Hits.Hits))
	}

	return searchResult.Hits.Hits[0].ID, nil
}

type post struct {
	ID      int    `json:"id"`
	Caption string `json:"caption"`
}

//these operator c, r
func createPost(post *post, client *elasticsearch.Client) error {

	instaPost := insta_post_text_filter.InstaPost{PostId: strconv.Itoa(post.ID), Caption: post.Caption}

	instaPostJson, err := json.Marshal(instaPost)

	if err != nil {
		panic(err)
	}
	response, err := client.Index("insta_post", strings.NewReader(string(instaPostJson)))

	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("FindPostId Search Query Failed StatusCode: %d", response.StatusCode)
	}

	return nil
}

//these operator u
func updateDocumentId(documentId string, caption string, client *elasticsearch.Client) error {
	instaPost := fmt.Sprintf(instaPostUpdate, caption)
	response, err := client.Update("insta_post", documentId, strings.NewReader(instaPost))

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("FindPostId Search Query Failed StatusCode: %d", response.StatusCode)
	}
	return nil
}
