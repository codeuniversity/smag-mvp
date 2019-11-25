package neo4jinserter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/codeuniversity/smag-mvp/worker"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"

	kf "github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	*worker.Worker

	qReader *kf.Reader
	conn    bolt.Conn

	inserterFunc InserterFunc
}

// InserterFunc is responsible to unmashal to the
// needed Data from the change Message and inserts
// it into neo4j
type InserterFunc func(*changestream.ChangeMessage, bolt.Conn) error

// New returns an initilized scraper
func New(neo4jConfig *utils.Neo4jConfig, userQReader *kf.Reader, inserterFunc InserterFunc) *Inserter {
	i := &Inserter{}

	i.qReader = userQReader
	i.inserterFunc = inserterFunc

	conn, err := initializeNeo4j(neo4jConfig)

	if err != nil {
		panic(err)
	}
	log.Println("✅ Neo4j Connection established")
	i.conn = conn

	i.Worker = worker.Builder{}.WithName("neo4j-inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("userQReader", userQReader.Close).
		AddShutdownHook("Neo4j", conn.Close).
		MustBuild()

	return i
}

// runStep the inserter
func (i *Inserter) runStep() error {
	m, err := i.qReader.FetchMessage(context.Background())

	if err != nil {
		return err
	}

	changeMessage := &changestream.ChangeMessage{}

	err = json.Unmarshal(m.Value, changeMessage)

	if err != nil {
		return err
	}

	err = i.inserterFunc(changeMessage, i.conn)

	if err != nil {
		return err
	}

	log.Println("Inserted")
	return i.qReader.CommitMessages(context.Background(), m)
}

//initializeNeo4j sets connection and constraints for neo4j
func initializeNeo4j(config *utils.Neo4jConfig) (bolt.Conn, error) {
	driver := bolt.NewDriver()
	address := fmt.Sprintf("bolt://%s:%s@%s:7687", config.Username, config.Password, config.Host)
	con, err := driver.OpenNeo(address)

	if err != nil {
		return nil, err
	}

	_, err = con.ExecNeo("CREATE CONSTRAINT ON (U:USER) ASSERT U.id IS UNIQUE", nil)

	if err != nil {
		return nil, err
	}

	return con, nil
}
