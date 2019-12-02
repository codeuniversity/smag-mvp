package main

import (
	"log"
	"os"
	"strconv"

	analyzer "github.com/codeuniversity/smag-mvp/nlp/frequency-analyzer"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Num of args isn't 2.")
	}
	username := os.Args[1]
	numOfWords, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("Unable to convert os.ARgs[2] to integer.")
	}

	keywordMap := analyzer.GetMostFrequentTerms(username, numOfWords)

	log.Printf("10 most frequently used words for user %s:\n%+v", username, keywordMap)
}
