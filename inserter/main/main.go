package main

import (
	"github.com/alexmorten/instascraper/inserter"
)

func main() {
	i := inserter.New()
	i.Run()
}
