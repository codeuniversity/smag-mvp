package kafka

import (
	"time"

	"github.com/codeuniversity/smag-mvp/utils"
	"github.com/segmentio/kafka-go"
)

// ReaderConfig is an internal helper structure for creating
// a new pre-configurated kafka-go Reader
type ReaderConfig struct {
	Address string
	GroupID string
	Topic   string
}

// WriterConfig is an internal helper structure for creating
// a new pre-configurated kafka-go Writer
type WriterConfig struct {
	Address string
	Topic   string
	Async   bool
}

// NewReaderConfig creates a new ReaderConfig structure
func NewReaderConfig(kafkaAddress, groupID, topic string) *ReaderConfig {
	return &ReaderConfig{
		Address: kafkaAddress,
		GroupID: groupID,
		Topic:   topic,
	}
}

// NewWriterConfig creates a new WriterConfig structure
func NewWriterConfig(kafkaAddress, topic string, async bool) *WriterConfig {
	return &WriterConfig{
		Address: kafkaAddress,
		Topic:   topic,
		Async:   async,
	}
}

// NewReader creates a kafka-go Reader structure by using common
// configuration and additionally applying the ReaderConfig on it
func NewReader(c *ReaderConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:               []string{c.Address},
		GroupID:               c.GroupID,
		Topic:                 c.Topic,
		MinBytes:              1e3,   // 1KB
		MaxBytes:              100e6, // 10MB
		QueueCapacity:         10000,
		CommitInterval:        time.Second,
		ReadBackoffMax:        time.Second * 5,
		WatchPartitionChanges: true,
		RetentionTime:         time.Hour * 24 * 30,
	})
}

// NewWriter creates a kafka-go Writer structure by using common
// configurations and additionally applying the WriterConfig on it
func NewWriter(c *WriterConfig) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{c.Address},
		Topic:    c.Topic,
		Balancer: &kafka.LeastBytes{},
		Async:    c.Async,
	})
}

//GetInserterConfig returns the Reader topics from kafka for Inserters
func GetInserterConfig() *ReaderConfig {
	var readerConfig *ReaderConfig

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")

	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	rTopic := utils.MustGetStringFromEnv("KAFKA_INFO_TOPIC")

	readerConfig = NewReaderConfig(kafkaAddress, groupID, rTopic)

	return readerConfig
}

// GetScraperConfig is a convenience function for gathering the necessary
// kafka configuration for all golang scrapers
func GetScraperConfig() (*ReaderConfig, *WriterConfig, *WriterConfig) {
	var nameReaderConfig *ReaderConfig
	var infoWriterConfig *WriterConfig
	var errWriterConfig *WriterConfig

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")

	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	nameTopic := utils.MustGetStringFromEnv("KAFKA_NAME_TOPIC")
	infoTopic := utils.MustGetStringFromEnv("KAFKA_INFO_TOPIC")
	errTopic := utils.MustGetStringFromEnv("KAFKA_ERR_TOPIC")

	nameReaderConfig = NewReaderConfig(kafkaAddress, groupID, nameTopic)
	infoWriterConfig = NewWriterConfig(kafkaAddress, infoTopic, true)
	errWriterConfig = NewWriterConfig(kafkaAddress, errTopic, false)

	return nameReaderConfig, infoWriterConfig, errWriterConfig
}

// GetInstaPostsScraperConfig is a convenience function for gathering the necessary
// kafka configuration for the insta posts golang scrapers
func GetInstaPostsScraperConfig() (*ReaderConfig, *WriterConfig, *WriterConfig) {
	var nameReaderConfig *ReaderConfig
	var infoWriterConfig *WriterConfig
	var errWriterConfig *WriterConfig

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")

	groupID := utils.MustGetStringFromEnv("KAFKA_GROUPID")
	nameTopic := utils.MustGetStringFromEnv("KAFKA_NAME_TOPIC")
	postsTopic := utils.MustGetStringFromEnv("KAFKA_INSTA_POSTS_TOPIC")
	errTopic := utils.MustGetStringFromEnv("KAFKA_ERR_TOPIC")

	nameReaderConfig = NewReaderConfig(kafkaAddress, groupID, nameTopic)
	infoWriterConfig = NewWriterConfig(kafkaAddress, postsTopic, true)
	errWriterConfig = NewWriterConfig(kafkaAddress, errTopic, false)

	return nameReaderConfig, infoWriterConfig, errWriterConfig
}
