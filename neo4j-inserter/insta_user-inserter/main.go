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

	changesTopic := utils.GetStringFromEnvWithDefault("KAFKA_CHANGE_TOPIC", "postgres.public.follows")
	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	groupID := utils.GetStringFromEnvWithDefault("KAFKA_GROUPID", "neo4j-user")
	qReaderConfig := kafka.NewReaderConfig(kafkaAddress, groupID, changesTopic)

	i := inserter.New(neo4jHost, neo4jUsername, neo4jPassword, kafka.NewReader(qReaderConfig), insertUsersAndFollowings)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

type Follow struct {
	FromID int `json:"from_id"`
	ToID   int `json:"to_id"`
}

func insertUsersAndFollowings(m *changestream.ChangeMessage, conn bolt.Conn) error {
	const createUsersAndRelationships = `
	MERGE(u1:USER{id: {fromID}})
	MERGE(u2:USER{id: {toID}})
	MERGE(u1)-[:FOLLOWS]->(u2)
	`
	f := &Follow{}
	err := json.Unmarshal(m.Payload.After, f)

	if err != nil {
		return err
	}

	_, err = conn.ExecNeo(createUsersAndRelationships, map[string]interface{}{"fromID": f.FromID, "toID": f.ToID})

	if err != nil {
		return err
	}

	return nil
}
