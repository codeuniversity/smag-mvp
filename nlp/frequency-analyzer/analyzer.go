package analyzer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/codeuniversity/smag-mvp/elastic"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

// Analyzer ...
type Analyzer struct {
	*elasticsearch.Client
}

// New ...
func New(esHosts []string) *Analyzer {
	return &Analyzer{
		Client: elastic.InitializeElasticSearch(esHosts),
	}
}

// MatchTermsForUser returns a list of identified terms
func (a *Analyzer) MatchTermsForUser(userID int, terms []string) ([]string, error) {
	foundTerms := make([]string, 0, len(terms))

	for _, term := range terms {
		didFindTerm, err := a.matchSingleTermForUser(userID, term)
		if err != nil {
			return nil, err
		}
		if didFindTerm == true {
			foundTerms = append(foundTerms, term)
		}
	}

	return foundTerms, nil
}

func (a *Analyzer) matchSingleTermForUser(userID int, term string) (bool, error) {
	searchBodyReader := newSearch(userID, term)
	esResponse, err := a.Search(
		a.Search.WithIndex("_all"),
		a.Search.WithBody(searchBodyReader),
	)
	if err != nil {
		return false, err
	}
	if esResponse.IsError() {
		return false, fmt.Errorf("Searching for term=%v failed status=%s body=%s", term, esResponse.Status(), esResponse.String())
	}
	body, err := ioutil.ReadAll(esResponse.Body)
	if err != nil {
		return false, err
	}

	r := &searchResponse{}
	if err = json.Unmarshal(body, r); err != nil {
		return false, err
	}

	if r.Hits.Total.Value > 0 {
		return true, nil
	}
	return false, nil
}

// nest is used to create a json-like map-compositions
type nest map[string]interface{}

func newSearch(userID int, term string) io.Reader {
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"term": map[string]interface{}{"user_id": userID}},
				},
				"should": []map[string]interface{}{
					{"match": map[string]interface{}{"caption": term}},
					{"match": map[string]interface{}{"comment": term}},
					{"match": map[string]interface{}{"bio": term}},
				}}},
	}

	return esutil.NewJSONReader(searchBody)
}

type searchResponse struct {
	Took     int          `json:"int"`
	TimedOut bool         `json:"timed_out"`
	Hits     responseHits `json:"hits"`
}

type responseHits struct {
	Total    total    `json:"total"`
	MaxScore float32  `json:"max_score"`
	Hits     []docHit `json:"hits"`
}

type total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type docHit struct {
	Index string  `json:"_index"`
	DocID string  `json:"_id"`
	Score float32 `json:"_score"`
	// Source hitDoc  `json:"_source"`
}

// type hitDoc struct {
// 	// post
// 	UserID  int    `json:"user_id"`
// 	Caption string `json:"caption"`
// 	// comment
// 	PostID  int    `json:"post_id"`
// 	Comment string `json:"comment_text"`
// 	// user
// 	Username string `json:"user_name"`
// 	Realname string `json:"real_name"`
// 	Bio      string `json:"bio"`
// }
