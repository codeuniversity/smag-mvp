package main

import (
	"encoding/json"
	"fmt"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	elasticsearch_inserter "github.com/codeuniversity/smag-mvp/kafka/changestream/elastic-inserter"
	insta_profile_filter "github.com/codeuniversity/smag-mvp/kafka/changestream/elastic-inserter/insta_profile-filter"
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
        "source": "ctx._source.bio = params.bio",
        "lang": "painless",
        "params" : {
            "bio" : "%s"
        }
    }
}
`

func main() {
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "my-kafka:9092")
	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.posts")

	elasticAddress := utils.MustGetStringFromEnv("ELASTIC_SEARCH_ADDRESS")

	elasticInserter := elasticsearch_inserter.New(elasticAddress, kafkaAddress, changesTopic, groupID, insertInstaProfile)

	service.CloseOnSignal(elasticInserter)
	waitUntilClosed := elasticInserter.Start()

	waitUntilClosed()
}

func insertInstaProfile(m *changestream.ChangeMessage, client *elasticsearch.Client) error {
	userProfile := &user{}
	err := json.Unmarshal(m.Payload.After, userProfile)

	if err != nil {
		return err
	}

	switch m.Payload.Op {
	case "r", "c":
		return createUserProfile(userProfile, client)
	case "u":
		documentId, err := findUserId(strconv.Itoa(userProfile.ID), client)

		if err != nil {
			return err
		}

		return updateDocumentId(documentId, userProfile.Bio, client)
	}

	return nil
}

//these operator u, o to find documentId
func findUserId(userId string, client *elasticsearch.Client) (string, error) {

	profileQuery := "userId: \"%s\""
	query := fmt.Sprintf(profileQuery, userId)
	response, err := client.Search(client.Search.WithIndex("insta_profile"), client.Search.WithQuery(query))

	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("FindPostId Search Query Failed StatusCode: %d", response.StatusCode)
	}

	searchResult := insta_profile_filter.QueryResult{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &searchResult)

	if err != nil {
		panic(err)
	}

	if len(searchResult.Hits.Hits) == 1 {
		return "", fmt.Errorf("found Duplicate Documents: UserId: %s Duplicates: %d", userId, len(searchResult.Hits.Hits))
	}

	return searchResult.Hits.Hits[0].ID, nil
}

//these operator c, r
func createUserProfile(userProfile *user, client *elasticsearch.Client) error {

	instaProfile := insta_profile_filter.InstaProfile{UserId: strconv.Itoa(userProfile.ID), Bio: userProfile.Bio}

	instaProfileJson, err := json.Marshal(instaProfile)

	if err != nil {
		panic(err)
	}
	response, err := client.Index("insta_profile", strings.NewReader(string(instaProfileJson)))

	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("createUserProfile Search Query Failed StatusCode: %d", response.StatusCode)
	}

	return nil
}

//these operator u
func updateDocumentId(documentId string, bio string, client *elasticsearch.Client) error {
	instaProfile := fmt.Sprintf(instaPostUpdate, bio)
	response, err := client.Update("insta_profile", documentId, strings.NewReader(instaProfile))

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("FindPostId Search Query Failed StatusCode: %d", response.StatusCode)
	}
	return nil
}

type user struct {
	ID  int    `json:"id"`
	Bio string `json:"bio"`
}
