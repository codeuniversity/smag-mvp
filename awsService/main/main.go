package main

import (
	"github.com/codeuniversity/smag-mvp/awsService"
)

func main() {
	s := awsService.New(9900)
	s.Listen()
}
