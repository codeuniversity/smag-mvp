package main

import (
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	elasticsearch_inserter "github.com/codeuniversity/smag-mvp/kafka/changestream/elastic-inserter"
	insta_comments_filter "github.com/codeuniversity/smag-mvp/kafka/changestream/elastic-inserter/insta_comments-filter"
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
        "source": "ctx._source.comment = params.comment",
        "lang": "painless",
        "params" : {
            "comment" : "%s"
        }
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
		return createComment(comment, client)
	case "u":
		documentId, err := findCommentId(strconv.Itoa(comment.ID), client)

		if err != nil {
			return err
		}

		return updateDocumentId(documentId, comment.Comment, client)
	}

	return nil
}

//these operator u, o to find documentId
func findCommentId(commentId string, client *elasticsearch.Client) (string, error) {

	postQuery := "commentId: \"%s\""
	query := fmt.Sprintf(postQuery, commentId)
	response, err := client.Search(client.Search.WithIndex("insta_comments"), client.Search.WithQuery(query))

	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("FindPostId Search Query Failed StatusCode: %d", response.StatusCode)
	}

	searchResult := insta_comments_filter.QueryResult{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &searchResult)

	if err != nil {
		panic(err)
	}

	if len(searchResult.Hits.Hits) == 1 {
		return "", fmt.Errorf("found Duplicate Documents: PostId: %s Duplicates: %d", commentId, len(searchResult.Hits.Hits))
	}

	return searchResult.Hits.Hits[0].ID, nil
}

type comment struct {
	ID      int    `json:"id"`
	PostId  string `json:"post_id"`
	Comment string `json:"comment_text"`
}

//these operator c, r
func createComment(comment *comment, client *elasticsearch.Client) error {

	instaPost := insta_comments_filter.InstaComment{CommentID: strconv.Itoa(comment.ID), PostID: comment.PostId, Comment: comment.Comment}

	instaPostJson, err := json.Marshal(instaPost)

	if err != nil {
		panic(err)
	}
	response, err := client.Index("insta_comments", strings.NewReader(string(instaPostJson)))

	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("createComment Create Post Failed StatusCode: %d", response.StatusCode)
	}

	return nil
}

//these operator u
func updateDocumentId(documentId string, comment string, client *elasticsearch.Client) error {
	instaPost := fmt.Sprintf(instaPostUpdate, comment)
	response, err := client.Update("insta_comments", documentId, strings.NewReader(instaPost))

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("FindPostId Update Document Failed StatusCode: %d", response.StatusCode)
	}
	return nil
}
