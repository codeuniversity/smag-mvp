package main

import (
	"encoding/json"
	"github.com/codeuniversity/smag-mvp/neo4j-inserter"

	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
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
