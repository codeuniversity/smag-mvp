package main

import neo4j_dump "github.com/codeuniversity/smag-mvp/neo4j-create-dump"

func main() {

	groupId := "test"
	topic := "inserterTopic"
	kafkaAddress := "my-kafka:9092"
	kafkaChunk := 100

	importDump := neo4j_dump.New(kafkaAddress, topic, groupId, kafkaChunk)

	importDump.Run()
}
