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

	kf "github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	*worker.Worker

	qReader *kf.Reader
	conn    bolt.Conn

	inserFunc InserterFunc
}

type InserterFunc func(*changestream.ChangeMessage, bolt.Conn) error

// New returns an initilized scraper
func New(neo4jAddress, neo4jUsername, neo4jPassword string, userQReader *kf.Reader, inserterFunc InserterFunc) *Inserter {
	i := &Inserter{}

	i.qReader = userQReader
	i.inserFunc = inserterFunc

	conn, err := initializeNeo4j(neo4jUsername, neo4jPassword, neo4jAddress)

	if err != nil {
		panic(err)
	}
	log.Println("✅ Neo4j Connection established")
	i.conn = conn

	i.Worker = worker.Builder{}.WithName("neo4j-inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("userQReader", userQReader.Close).
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

	err = i.inserFunc(changeMessage, i.conn)

	if err != nil {
		return err
	}

	log.Println("Inserted")
	return i.qReader.CommitMessages(context.Background(), m)
}

// insertUser creates user in neo4j if thats not the case already

//initializeNeo4j sets connection and constraints for neo4j
func initializeNeo4j(neo4jUsername, neo4jPassword, neo4jAddress string) (bolt.Conn, error) {
	driver := bolt.NewDriver()
	address := fmt.Sprintf("bolt://%s:%s@%s:7687", neo4jUsername, neo4jPassword, neo4jAddress)
	fmt.Println(neo4jPassword)
	con, err := driver.OpenNeo(address)

	if err != nil {
		return nil, err
	}

	// _, err = con.ExecNeo("CREATE CONSTRAINT ON (U:USER) ASSERT U.id IS UNIQUE", nil)

	if err != nil {
		return nil, err
	}

	return con, nil
}
