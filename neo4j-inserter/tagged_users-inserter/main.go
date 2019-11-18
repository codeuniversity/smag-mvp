package main

import (
	"encoding/json"

	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	inserter "github.com/codeuniversity/smag-mvp/neo4j-inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func main() {
	neo4jHost := utils.GetStringFromEnvWithDefault("NEO4J_HOST", "localhost")
	neo4jUsername := utils.GetStringFromEnvWithDefault("NEO4J_USERNAME", "neo4j")
	neo4jPassword := utils.GetStringFromEnvWithDefault("NEO4J_PASSWORD", "12345678")

	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.post_tagged_users")
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.GetStringFromEnvWithDefault("KAFKA_GROUPID", "neo4j-tagged-users")
	qReaderConfig := kafka.NewReaderConfig(kafkaAddress, groupID, changesTopic)

	i := inserter.New(neo4jHost, neo4jUsername, neo4jPassword, kafka.NewReader(qReaderConfig), addTaggedUsersRelationship)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

type taggedUser struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_	id"`
}

func addTaggedUsersRelationship(m *changestream.ChangeMessage, conn bolt.Conn) error {
	const addTaggedRelationship = `
	MERGE(u:USER{id: {userID}})
	MERGE(p:POST{id: {postID}})
	MERGE(u)-[:TAGGED_ON]->(p)
	`
	t := &taggedUser{}
	err := json.Unmarshal(m.Payload.After, t)

	if err != nil {
		return err
	}

	_, err = conn.ExecNeo(addTaggedRelationship, map[string]interface{}{"userID": t.UserID, "postID": t.PostID})

	if err != nil {
		return err
	}

	return nil
}
