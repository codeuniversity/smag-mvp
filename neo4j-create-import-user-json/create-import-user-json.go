package neo4j_import

import (
	"context"
	"encoding/json"
	"fmt"
	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/codeuniversity/smag-mvp/worker"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type Neo4jImport struct {
	kReader        *kafka.Reader
	kafkaChunkSize int
	*worker.Worker
	file *os.File
}

// New returns an initialised Indexer
func New(kafkaAddress, changesTopic, kafkaGroupID string, kafkaChunkSize int) *Neo4jImport {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)

	i := &Neo4jImport{}
	i.kReader = kf.NewReader(readerConfig)
	i.kafkaChunkSize = kafkaChunkSize

	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("neo4j-json-data-%s.json", id.String())
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	i.file = file

	if _, err = i.file.WriteString(startJson); err != nil {
		panic(err)
	}

	i.Worker = worker.Builder{}.WithName("insta_comments_scraper").
		WithWorkStep(i.runStep).
		AddShutdownHook("kReader", i.kReader.Close).
		AddShutdownHook("writeTheEndJson", i.writeTheEndJson).
		MustBuild()

	return i
}

type Follow struct {
	FromID int `json:"from_id"`
	ToID   int `json:"to_id"`
}

const startJson = `
	{
  "followers": [
`

const endJson = `
]
}
`

var isFirstWrite = true

func (i *Neo4jImport) runStep() error {

	messages, err := i.readMessageBlock(10*time.Second, i.kafkaChunkSize)

	log.Println("Messages Bulk: ", len(messages))
	if len(messages) == 0 {
		return nil
	}

	if err != nil {
		panic(err)
	}

	follows := make([]Follow, i.kafkaChunkSize)
	for _, message := range messages {

		changeMessage := &changestream.ChangeMessage{}
		if err := json.Unmarshal(message.Value, changeMessage); err != nil {
			panic(err)
		}

		log.Println("Change Message: ", string(changeMessage.Payload.After))

		f := &Follow{}
		err := json.Unmarshal(changeMessage.Payload.After, f)

		if err != nil {
			return err
		}

		follows = append(follows, *f)
	}

	var followsJson string
	for i, follow := range follows {
		followJson, err := json.Marshal(follow)

		if err != nil {
			return err
		}

		if !isFirstWrite || i != 0 {
			followsJson += ","
		}
		followsJson += string(followJson)
	}

	if _, err = i.file.WriteString(followsJson); err != nil {
		return err
	}
	isFirstWrite = false

	err = i.kReader.CommitMessages(context.Background(), messages...)
	if err != nil {
		return err
	}

	return nil
}

func (i *Neo4jImport) writeTheEndJson() error {
	if _, err := i.file.WriteString(endJson); err != nil {
		return err
	}
	return nil
}

func (i *Neo4jImport) readMessageBlock(timeout time.Duration, maxChunkSize int) (messages []kafka.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for k := 0; k < maxChunkSize; k++ {
		m, err := i.kReader.FetchMessage(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				return messages, nil
			}

			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}
