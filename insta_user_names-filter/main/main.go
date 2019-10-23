package main

import (
	filter "github.com/codeuniversity/smag-mvp/insta_user_names-filter"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	var t *filter.Transferer

	t = filter.New("postgres.public.users", "user_names")

	service.CloseOnSignal(t)
	waitUntilClosed := t.Start()

	waitUntilClosed()
}
