package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type ReaderConfig struct {
	Address string
	GroupID string
	Topic   string
}

type WriterConfig struct {
	Address string
	Topic   string
	Async   bool
}

func NewReaderConfig(kafkaAddress, groupID, topic string) *ReaderConfig {
	return &ReaderConfig{
		Address: kafkaAddress,
		GroupID: groupID,
		Topic:   topic,
	}
}

func NewWriterConfig(kafkaAddress, topic string, async bool) *WriterConfig {
	return &WriterConfig{
		Address: kafkaAddress,
		Topic:   topic,
		Async:   async,
	}
}

func NewReader(c *ReaderConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{c.Address},
		GroupID:        c.GroupID,
		Topic:          c.Topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
}

func NewWriter(c *WriterConfig) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{c.Address},
		Topic:    c.Topic,
		Balancer: &kafka.LeastBytes{},
		Async:    c.Async,
	})
}
