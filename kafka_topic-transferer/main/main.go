package main

import (
	transferer "github.com/codeuniversity/smag-mvp/kafka_topic-transferer"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	var t *transferer.Transferer

	t = transferer.New("postgres.public.users", "user_names")

	service.CloseOnSignal(t)
	go t.Run()

	t.WaitUntilClosed()
}
