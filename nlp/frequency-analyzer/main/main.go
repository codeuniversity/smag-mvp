package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	analyzer "github.com/codeuniversity/smag-mvp/nlp/frequency-analyzer"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Num of args isn't 1.\nUsage: go run main.go <user_id:int>")
	}
	userID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to convert user_id=%+v to integer.\nsage: go run main.go <user_id:int>", userID)
	}

	a := analyzer.New([]string{"http://127.0.0.1:9200"})
	log.Printf("a=%+v", a)

	// TODO: load cities.json
	jsonFile, err := os.Open("nlp/frequency-analyzer/cities.json")
	defer jsonFile.Close()
	if err != nil {
		log.Fatalln(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	var cityMap map[string][]string
	if err := json.Unmarshal(byteValue, &cityMap); err != nil {
		panic(err)
	}
	log.Printf("cityMap=%+v", cityMap)

	foundCities := make(map[string]bool)
	for city, cityTerms := range cityMap {
		foundTerms, err := a.MatchTermsForUser(userID, cityTerms)
		if err != nil {
			panic(err)
		}
		log.Printf("city=%v \t-> foundTerms=%+v", city, foundTerms)
		// check if there are results for city
		if len(foundTerms) > 0 {
			foundCities[city] = true
		} else {
			foundCities[city] = false
		}
	}
	log.Printf("foundCities=%+v", foundCities)

	res := make([]string, 0, len(foundCities))
	log.Printf("Could identify following cities for user=%v: {", userID)
	for city, found := range foundCities {
		if found == true {
			res = append(res, city)
			fmt.Printf("  * %v\n", city)
		}
	}
	fmt.Println("}")
}
