package main

import (
	"sync"

	"github.com/alexmorten/instascraper/inserter"
	"github.com/alexmorten/instascraper/scraper"
)

func main() {
	s := scraper.New()
	i := inserter.New()
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		s.Run()
		wg.Done()
	}()
	go func() {
		i.Run()
		wg.Done()
	}()

	wg.Wait()
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
