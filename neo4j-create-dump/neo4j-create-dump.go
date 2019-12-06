package neo4j_dump

import "C"
import (
	"context"
	"encoding/json"
	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/kafka/changestream"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type Neo4jImport struct {
	kReader        *kafka.Reader
	kafkaChunkSize int

	file *os.File
}

// New returns an initialised Indexer
func New(kafkaAddress, changesTopic, kafkaGroupID string, kafkaChunkSize int) *Neo4jImport {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)

	i := &Neo4jImport{}
	i.kReader = kf.NewReader(readerConfig)
	i.kafkaChunkSize = kafkaChunkSize

	file, err := os.OpenFile("neo4j-dump.json", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	i.file = file

	if _, err = i.file.WriteString(startJson); err != nil {
		panic(err)
	}
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

func (i *Neo4jImport) Run() {

	isFirstWrite := true
	for k := 0; k < 5; k++ {
		messages, err := i.readMessageBlock(10*time.Second, i.kafkaChunkSize)

		log.Println("Messages Bulk: ", len(messages))
		if len(messages) == 0 {
			continue
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
				panic(err)
			}

			follows = append(follows, *f)
		}

		var followsJson string
		for i, follow := range follows {
			followJson, err := json.Marshal(follow)

			if err != nil {
				panic(err)
			}

			if !isFirstWrite && i == 0 {
				followsJson += ","
			}
			followsJson += string(followJson)
		}

		if _, err = i.file.WriteString(followsJson); err != nil {
			panic(err)
		}
		isFirstWrite = false

		err = i.kReader.CommitMessages(context.Background(), messages...)
		if err != nil {
			panic(err)
		}
	}

	if _, err := i.file.WriteString(endJson); err != nil {
		panic(err)
	}

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
