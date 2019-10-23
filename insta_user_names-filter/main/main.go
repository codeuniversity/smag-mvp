package main

import (
	filter "github.com/codeuniversity/smag-mvp/insta_user_names-filter"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	f := filter.New("postgres.public.users", "user_names")

	service.CloseOnSignal(f)
	waitUntilClosed := f.Start()

	waitUntilClosed()
}
