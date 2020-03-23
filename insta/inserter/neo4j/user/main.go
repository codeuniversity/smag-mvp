package main

import (
	"encoding/json"

	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	neo4jinserter "github.com/codeuniversity/smag-mvp/neo4j/inserter"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {
	readerConfig := kafka.GetInserterConfig()
	neo4jConfig := utils.GetNeo4jConfig()

	i := neo4jinserter.New(neo4jConfig, kafka.NewReader(readerConfig), insertUsersAndFollowings)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

type Follow struct {
	FromID int `json:"from_id"`
	ToID   int `json:"to_id"`
}

func insertUsersAndFollowings(m *changestream.ChangeMessage, session neo4j.Session) error {
	const createUsersAndRelationships = `
	MERGE(u1:USER{id: $fromID})
	MERGE(u2:USER{id: $toID})
	MERGE(u1)-[:FOLLOWS]->(u2)
	`
	f := &Follow{}
	err := json.Unmarshal(m.Payload.After, f)

	if err != nil {
		return err
	}

	_, err = session.Run(createUsersAndRelationships, map[string]interface{}{"fromID": f.FromID, "toID": f.ToID})

	if err != nil {
		return err
	}

	return nil
}
