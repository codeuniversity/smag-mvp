package main

import (
	"log"
	"os"
	"strconv"

	analyzer "github.com/codeuniversity/smag-mvp/nlp/frequency-analyzer"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Num of args isn't 2.\nUsage: go run main.go <user_id:int>")
	}
	userID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to convert user_id=%+v to integer.\nsage: go run main.go <user_id:int>", userID)
	}

	a := analyzer.New([]string{"localhost:9200"})

	// TODO: load cities.json
	jsonFile, err := os.Open("cities.json")
	defer jsonFile.Close()
	if err != nil {
		log.Fatalln(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	var cityMap map[string][]string
	if err := json.Unmarshal(byteValue, &cityMap); err != nil {
		panic(err)
	}

	foundCities := make(map[string]bool)
	for city, cityTerms := range cityMap {
		foundTerms, err := a.MatchTermsForUser(userID, cityTerms)
		if err != nil {
			panic(err)
		}
		// check if there are results for city
		if len(foundTerms) > 0 {
			foundCities[city] = true
		} else {
			foundCities[city] = false
		}
	}

	log.Printf("Could identify following cities for user=%v:\n{%v}", userID, foundCities)
}
