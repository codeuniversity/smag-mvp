package main

import (
	"strconv"

	extractor "github.com/codeuniversity/smag-mvp/dgraph-extractor"
	"github.com/codeuniversity/smag-mvp/service"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	dgraphAddress := utils.GetStringFromEnvWithDefault("DGRPAH_ADDRESS", "127.0.0.1:9080")
	startID := utils.GetStringFromEnvWithDefault("START_ID", "1")

	startIDInt, err := strconv.ParseInt(startID, 10, 64)
	utils.PanicIfErr(err)
	i := extractor.New(kafkaAddress, dgraphAddress, int(startIDInt))

	service.CloseOnSignal(i)
	go i.Run()

	i.WaitUntilClosed()
}
