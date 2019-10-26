package main

import (
	"github.com/codeuniversity/smag-mvp/aws_service"
)

func main() {
	s := aws_service.New(9900)
	s.Listen()
}
