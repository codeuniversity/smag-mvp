package main

import (
	"log"
	"os"
	"strconv"

	analyzer "github.com/codeuniversity/smag-mvp/nlp/frequency-analyzer"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Num of args isn't 2.\nUsage: go run main.go <user_name:string> <num_of_words:int>")
	}
	username := os.Args[1]
	numOfWords, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("Unable to convert os.Args[2] to integer.\nUsage: go run main.go <user_name:string> <num_of_words:int>")
	}

	keywordMap := analyzer.GetMostFrequentTerms(username, numOfWords)

	log.Printf("- %v most frequently used words for user %s: { %v }", numOfWords, username, keywordMap)
}
