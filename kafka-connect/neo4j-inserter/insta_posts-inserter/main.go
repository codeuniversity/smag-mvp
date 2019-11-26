package main

import (
	"encoding/json"
	"github.com/codeuniversity/smag-mvp/kafka-connect/neo4j-inserter"

	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func main() {
	readerConfig := kafka.GetInserterConfig()
	neo4jConfig := utils.GetNeo4jConfig()

	i := neo4jinserter.New(neo4jConfig, kafka.NewReader(readerConfig), insertPostsAndAddRelationship)

	service.CloseOnSignal(i)
	waitUntilClosed := i.Start()

	waitUntilClosed()
}

type Post struct {
	UserID int `json:"user_id"`
	PostID int `json:"id"`
}

func insertPostsAndAddRelationship(m *changestream.ChangeMessage, conn bolt.Conn) error {
	const insertPostsAndAddRelationship = `
	MERGE(u:USER{id: {userID}})
	MERGE(p:POST{id: {postID}})
	MERGE(u)-[:POSTED]->(p)
	`
	p := &Post{}
	err := json.Unmarshal(m.Payload.After, p)

	if err != nil {
		return err
	}

	_, err = conn.ExecNeo(insertPostsAndAddRelationship, map[string]interface{}{"userID": p.UserID, "postID": p.PostID})

	if err != nil {
		return err
	}

	return nil
}
