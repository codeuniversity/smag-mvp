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

	"github.com/neo4j/neo4j-go-driver/neo4j"
	kf "github.com/segmentio/kafka-go"
)

// Inserter represents the scraper containing all clients it uses
type Inserter struct {
	*worker.Worker

	qReader *kf.Reader
	driver  neo4j.Driver
	session neo4j.Session

	inserterFunc InserterFunc
}

// InserterFunc is responsible to unmashal to the
// needed Data from the change Message and inserts
// it into neo4j
type InserterFunc func(*changestream.ChangeMessage, neo4j.Session) error

// New returns an initilized scraper
func New(neo4jConfig *utils.Neo4jConfig, userQReader *kf.Reader, inserterFunc InserterFunc) *Inserter {
	i := &Inserter{}

	i.qReader = userQReader
	i.inserterFunc = inserterFunc

	session, driver, err := initializeNeo4j(neo4jConfig)
	if err != nil {
		panic(err)
	}
	i.session = session
	i.driver = driver

	log.Println("✅ Neo4j Connection established")

	i.Worker = worker.Builder{}.WithName("neo4j-inserter").
		WithWorkStep(i.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("userQReader", userQReader.Close).
		AddShutdownHook("Neo4j Driver", driver.Close).
		AddShutdownHook("Neo4j Session", session.Close).
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

	err = i.inserterFunc(changeMessage, i.session)

	if err != nil {
		return err
	}

	log.Println("Inserted")
	return i.qReader.CommitMessages(context.Background(), m)
}

//initializeNeo4j sets connection and constraints for neo4j
func initializeNeo4j(config *utils.Neo4jConfig) (neo4j.Session, neo4j.Driver, error) {
	address := fmt.Sprintf("bolt://%s:7687", config.Host)
	driver, err := neo4j.NewDriver(address, neo4j.BasicAuth(config.Username, config.Password, ""))
	if err != nil {
		return nil, nil, err
	}

	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, nil, err
	}

	_, err = session.Run("CREATE CONSTRAINT ON (U:USER) ASSERT U.id IS UNIQUE", nil)
	if err != nil {
		return nil, nil, err
	}

	return session, driver, nil
}
