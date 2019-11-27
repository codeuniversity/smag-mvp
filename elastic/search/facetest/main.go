package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codeuniversity/smag-mvp/elastic/search/faces"
	"github.com/codeuniversity/smag-mvp/faces/proto"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"google.golang.org/grpc"
)

func main() {

	if len(os.Args) != 2 {
		panic("requires exactly one param - the url to an image")
	}

	con, err := grpc.Dial("localhost:6666", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	faceReconClient := proto.NewFaceRecognizerClient(con)
	cfg := elasticsearch.Config{Addresses: []string{"http://localhost:9200"}}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating elasticsearch client: %s", err)
	}

	client := &faces.FaceSearchClient{
		FaceRecognitionClient: faceReconClient,
		ESClient:              esClient,
	}

	faces, err := client.FindSimilarFacesInImage(os.Args[1], 10)

	fmt.Printf("\n (")
	for _, face := range faces {
		fmt.Printf("%d,", face.PostID)
	}
	fmt.Printf(")\n")
}
