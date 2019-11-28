package faces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codeuniversity/smag-mvp/elastic"
	proto "github.com/codeuniversity/smag-mvp/faces/proto"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"io/ioutil"
	"log"
)

// Client searches elasticsearch for similar faces to the faces given in an image
type Client struct {
	FaceRecognitionClient proto.FaceRecognizerClient
	ESClient              *elasticsearch.Client
}

// FindSimilarFacesInImage finds faces in the given image and searches for faces similar to the ones in the given image
// the url should be downloadable by the face recognizer, i.e. a signed url for s3 or for the imgproxy
func (c *Client) FindSimilarFacesInImage(imgURL string, maxHitsPerFace int) ([]FoundFace, error) {
	response, err := c.FaceRecognitionClient.RecognizeFaces(context.Background(), &proto.RecognizeRequest{Url: imgURL})
	if err != nil {
		return nil, err
	}
	faces := response.Faces
	log.Println("found ", len(faces), " faces in given image")
	foundFaces := []FoundFace{}
	for _, face := range faces {
		r, err := c.searchForSimilarFaces(face, maxHitsPerFace)
		if err != nil {
			return nil, err
		}
		log.Println("max score:", r.Hits.MaxScore)
		for _, hit := range r.Hits.Hits {
			doc := hit.Source
			log.Println("score:", hit.Score)
			foundFaces = append(foundFaces, FoundFace{FaceDoc: doc, Score: hit.Score, MaxScoreShare: hit.Score / r.Hits.MaxScore})
		}
	}

	return foundFaces, nil
}

func (c *Client) searchForSimilarFaces(face *proto.Face, maxHits int) (*searchResponse, error) {
	searchBodyReader := esutil.NewJSONReader(newSearch(maxHits, face.Encoding))
	esResponse, err := c.ESClient.Search(
		c.ESClient.Search.WithIndex(elastic.FacesIndex),
		c.ESClient.Search.WithBody(searchBodyReader),
	)
	if err != nil {
		return nil, err
	}
	if esResponse.IsError() {
		return nil, fmt.Errorf("searching for face failed status=%s body=%s", esResponse.Status(), esResponse.String())
	}
	body, err := ioutil.ReadAll(esResponse.Body)
	if err != nil {
		return nil, err
	}

	response := &searchResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type searchBody struct {
	Size  int         `json:"size"`
	Query searchQuery `json:"query"`
}

type searchQuery struct {
	FunctionScore searchFunctionScore `json:"function_score"`
}

type searchFunctionScore struct {
	BoostMode   string            `json:"boost_mode"`
	ScriptScore searchScriptScore `json:"script_score"`
}

type searchScriptScore struct {
	Script searchScript `json:"script"`
}

type searchScript struct {
	Source string       `json:"source"`
	Lang   string       `json:"lang"`
	Params searchParams `json:"params"`
}

type searchParams struct {
	Cosine bool      `json:"cosine"`
	Field  string    `json:"field"`
	Vector []float32 `json:"vector"`
}

type searchResponse struct {
	Took     int          `json:"int"`
	TimedOut bool         `json:"timed_out"`
	Hits     responseHits `json:"hits"`
}

type responseHits struct {
	MaxScore float32  `json:"max_score"`
	Hits     []docHit `json:"hits"`
}

type docHit struct {
	Index  string  `json:"_index"`
	DocID  string  `json:"_id"`
	Score  float32 `json:"_score"`
	Source FaceDoc `json:"_source"`
}

// FaceDoc is the face hit that elasticsearch returned for the search
type FaceDoc struct {
	PostID   int    `json:"post_id"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"Width"`
	Height   int    `json:"Height"`
	Encoding string `json:"encoding_vector"`
}

// FoundFace includes scoring information
type FoundFace struct {
	FaceDoc
	Score         float32
	MaxScoreShare float32
}

func newSearch(size int, encoding []float32) *searchBody {
	return &searchBody{
		Size: size,
		Query: searchQuery{
			FunctionScore: searchFunctionScore{
				BoostMode: "replace",
				ScriptScore: searchScriptScore{
					Script: searchScript{
						Source: "binary_vector_score",
						Lang:   "knn",
						Params: searchParams{
							Cosine: true,
							Field:  "encoding_vector",
							Vector: encoding,
						},
					},
				},
			},
		},
	}
}
