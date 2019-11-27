package main

import (
	"fmt"
	"os"

	"github.com/codeuniversity/smag-mvp/elastic"
	"github.com/codeuniversity/smag-mvp/elastic/search/faces"
	"github.com/codeuniversity/smag-mvp/faces/proto"
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
	esClient := elastic.InitializeElasticSearch([]string{"http://localhost:9200"})

	client := &faces.Client{
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
