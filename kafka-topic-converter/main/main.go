package main

import (
	"github.com/codeuniversity/smag-mvp/kafka"
	conv "github.com/codeuniversity/smag-mvp/kafka-topic-transferer"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	converter := &conv.Converter{}

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	readerConfig := kafka.NewReaderConfig(kafkaAddress, "22", "postgres.public.users")
	writerConfig := kafka.NewWriterConfig(kafkaAddress, "user_names", true)
	reader := kafka.NewReader(readerConfig)
	writer := kafka.NewWriter(writerConfig)

	converter = conv.New(reader, writer)

	converter.Run()
}
