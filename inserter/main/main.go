package main

import (
	"github.com/alexmorten/instascraper/inserter"
	"github.com/alexmorten/instascraper/utils"
)

func main() {
	i := inserter.New()

	utils.CloseOnSignal(i.Close)
	go i.Run()

	i.WaitUntilClosed()
}
