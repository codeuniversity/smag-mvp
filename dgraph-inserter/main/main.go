package main

import (
	inserter "github.com/codeuniversity/smag-mvp/dgraph-inserter"
	"github.com/codeuniversity/smag-mvp/kafka"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {
	var i *inserter.Inserter

	dgraphAddress := utils.GetStringFromEnvWithDefault("DGRPAH_ADDRESS", "127.0.0.1:9080")

	isUserDiscovery, err := utils.GetBoolFromEnvWithDefault("USER_DISCOVERY", false)
	if err != nil {
		panic(err)
	}

	qReaderConfig, qWriterConfig := kafka.GetInserterConfig(isUserDiscovery)

	if isUserDiscovery {
		i = inserter.New(
			dgraphAddress,
			kafka.NewReader(qReaderConfig),
			kafka.NewWriter(qWriterConfig),
		)
	} else {
		i = inserter.New(
			dgraphAddress,
			kafka.NewReader(qReaderConfig),
			nil,
		)
	}

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
