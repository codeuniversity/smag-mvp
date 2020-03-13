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

	i := neo4jinserter.New(neo4jConfig, kafka.NewReader(readerConfig), addTaggedUsersRelationship)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

type taggedUser struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

func addTaggedUsersRelationship(m *changestream.ChangeMessage, session neo4j.Session) error {
	const addTaggedRelationship = `
	MERGE(u:USER{id: $userID})
	MERGE(p:POST{id: $postID})
	MERGE(u)-[:TAGGED_ON]->(p)
	`
	t := &taggedUser{}
	err := json.Unmarshal(m.Payload.After, t)

	if err != nil {
		return err
	}

	_, err = session.Run(addTaggedRelationship, map[string]interface{}{"userID": t.UserID, "postID": t.PostID})

	if err != nil {
		return err
	}

	return nil
}
