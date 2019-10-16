package main

import (
	conv "github.com/codeuniversity/smag-mvp/kafka_topic_transferer"
)

func main() {
	converter := &conv.Converter{}

	converter = conv.New("postgres.public.users", "user_names")

	converter.Run()
}
