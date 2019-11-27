package main

import (
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/elastic"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"strconv"

	esIndexer "github.com/codeuniversity/smag-mvp/elastic/indexer"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.GetStringFromEnvWithDefault("KAFKA_GROUPID", "insta_usersearch-inserter")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.users")

	esHosts := utils.GetMultipliesStringsFromEnvDefault("ELASTIC_SEARCH_ADDRESS", []string{"http://localhost:9201"})

	// create and run esInserter
	i := esIndexer.New(esHosts, elastic.UsersIndex, elastic.UsersIndexMapping, kafkaAddress, changesTopic, groupID, handleChangemessage)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

// handleChangemessage filters relevant events and upserts them
func handleChangemessage(esClient *elasticsearch.Client, m *changestream.ChangeMessage) error {
	user := &user{}
	if err := json.Unmarshal(m.Payload.After, user); err != nil {
		return err
	}

	switch m.Payload.Op {
	case "c", "r", "u":
		return upsertDocument(user, esClient)
	}

	return nil
}

func upsertDocument(u *user, esClient *elasticsearch.Client) error {
	upsertBody := createUpsertBody(u)
	response, err := esClient.Update(
		elastic.UsersIndex,
		strconv.Itoa(u.ID),
		esutil.NewJSONReader(upsertBody))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to upsert user. StatusCode: %d", response.StatusCode)
	}

	return nil
}

func createUpsertBody(user *user) map[string]interface{} {
	var commentUpsert = map[string]interface{}{
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

	return commentUpsert
}
