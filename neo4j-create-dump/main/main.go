package main

import neo4j_dump "github.com/codeuniversity/smag-mvp/neo4j-create-dump"

func main() {

	groupId := "neo4j-user-inserter4"
	topic := "postgres.public.follows"
	kafkaAddress := "my-kafka:9092"
	kafkaChunk := 10

	importDump := neo4j_dump.New(kafkaAddress, topic, groupId, kafkaChunk)

	importDump.Run()
}
