package neo4jinserter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/worker"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"

	"github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	*worker.Worker

	qReader *kafka.Reader
	conn    bolt.Conn
}

type user struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
}

// New returns an initilized scraper
func New(neo4jAddress, neo4jUsername, neo4jPassword string, userQReader *kafka.Reader) *Inserter {
	i := &Inserter{}

	conn, err := initializeNeo4j(neo4jUsername, neo4jPassword, neo4jAddress)
	if err != nil {
		panic(err)
	}
	i.conn = conn

	i.Worker = worker.Builder{}.WithName("neo4j_user-inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("userQReader", userQReader.Close).
		MustBuild()

	return i
}

// runStep the inserter
func (i *Inserter) runStep() error {

	log.Println("starting inserter")
	m, err := i.qReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	changeMessage := &changestream.ChangeMessage{}

	err = json.Unmarshal(m.Value, changeMessage)
	if err != nil {
		return err
	}

	err = i.insertUser(changeMessage)
	if err != nil {
		return err
	}

	log.Println("User inserted:")
	return i.qReader.CommitMessages(context.Background(), m)
}

// insertUser creates user in neo4j if thats not the case already
func (i *Inserter) insertUser(m *changestream.ChangeMessage) error {
	const createUser = `
		MERGE (u:USER {id: {id}, user_name: {userName}})
	`
	u := &user{}
	err := json.Unmarshal(m.Payload.After, u)
	if err != nil {
		return err
	}

	_, err = i.conn.ExecNeo(createUser, map[string]interface{}{"id": u.ID, "userName": u.UserName})
	if err != nil {
		return err
	}

	return nil
}

// sets connection and constraints for neo4j
func initializeNeo4j(neo4jUsername, neo4jPassword, neo4jAddress string) (bolt.Conn, error) {
	driver := bolt.NewDriver()
	address := fmt.Sprintf("bolt://%s:%s@%s", neo4jUsername, neo4jPassword, neo4jAddress)
	con, err := driver.OpenNeo(address)
	if err != nil {
		return nil, err
	}

	_, err = con.ExecNeo("CREATE CONSTRAINT ON (U:USER) ASSERT U.name IS UNIQUE", nil)
	if err != nil {
		return nil, err
	}

	return con, nil
}
