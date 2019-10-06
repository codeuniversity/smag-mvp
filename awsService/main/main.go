package main

import (
	"github.com/codeuniversity/smag-mvp/awsService"
	"github.com/codeuniversity/smag-mvp/service"
)

func main() {
	kafkaAddress := ""
	if kafkaAddress == "" {
		kafkaAddress = "52.58.171.160:9092"
	}
	s := awsService.New(kafkaAddress)
	service.CloseOnSignal(s)
	go s.Run()

	s.WaitUntilClosed()
}
