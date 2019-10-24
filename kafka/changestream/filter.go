package changestream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kf "github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/worker"

	"github.com/segmentio/kafka-go"
)

// Filter is responsible for reading the change stream,
// filtering out the events that are not interesting to us
// and writing new messages based on the changes to the filtered topic
type Filter struct {
	*worker.Worker

	changesReader  *kafka.Reader
	filteredWriter *kafka.Writer

	filterFunc FilterFunc
}

// FilterFunc given a ChangeMessage from the changesTopic
// returns zero, one or multiple kafka Messages that should be written to the filteredTopic
type FilterFunc func(*ChangeMessage) ([]kafka.Message, error)

// NewFilter returns an initilized Filter
func NewFilter(kafkaAddress, kafkaGroupID, changesTopic, filteredTopic string, filter FilterFunc) *Filter {
	readerConfig := kf.NewReaderConfig(kafkaAddress, kafkaGroupID, changesTopic)
	writerConfig := kf.NewWriterConfig(kafkaAddress, filteredTopic, true)

	f := &Filter{
		changesReader:  kf.NewReader(readerConfig),
		filteredWriter: kf.NewWriter(writerConfig),
		filterFunc:     filter,
	}

	b := worker.Builder{}.WithName(fmt.Sprintf("changestream_filter[%s->%s]", changesTopic, filteredTopic)).
		WithWorkStep(f.runStep).
		WithStopTimeout(10*time.Second).
		AddShutdownHook("changesReader", f.changesReader.Close).
		AddShutdownHook("filteredWriter", f.filteredWriter.Close)

	f.Worker = b.MustBuild()

	return f
}

func (t *Filter) runStep() error {
	m, err := t.changesReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	changeMessage := &ChangeMessage{}
	err = json.Unmarshal(m.Value, changeMessage)
	if err != nil {
		return err
	}

	kafkaMessages, err := t.filterFunc(changeMessage)
	if err != nil {
		return err
	}

	if len(kafkaMessages) > 0 {
		err = t.filteredWriter.WriteMessages(context.Background(), kafkaMessages...)
		if err != nil {
			return err
		}
	}

	return t.changesReader.CommitMessages(context.Background(), m)

}
