package generator

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// HTTPHeaderGenerator generates headers for http requests for scraping
type HTTPHeaderGenerator struct {
	browserAgent browserAgent
}

// New returns an initialized HTTPHeaderGenerator
func New() *HTTPHeaderGenerator {
	generator := &HTTPHeaderGenerator{}
	data, err := ioutil.ReadFile("http_header-generator/useragents.json")
	if err != nil {
		panic(err)
	}
	var userAgent browserAgent
	err = json.Unmarshal(data, &userAgent)

	if err != nil {
		panic(err)
	}
	generator.browserAgent = userAgent
	return generator
}

type browserAgent []struct {
	UserAgents string `json:"useragent"`
}

// AddHeaders ads the generated headers to the request headers
func (h *HTTPHeaderGenerator) AddHeaders(header *http.Header) {
	header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	header.Add("Accept-Charset", "utf-8")
	header.Add("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	header.Add("Cache-Control", "no-cache")
	header.Add("Content-Type", "application/json; charset=utf-8")
	header.Add("User-Agent", h.GetRandomUserAgent())
}

// GetRandomUserAgent returns a random user agent
func (h *HTTPHeaderGenerator) GetRandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(len(h.browserAgent))
	return h.browserAgent[randomNumber].UserAgents
}
