package main

import (
	transferer "github.com/codeuniversity/smag-mvp/kafka_topic-transferer"
)

func main() {

	t := transferer.New("postgres.public.users", "user_names")

	t.Run()
}
