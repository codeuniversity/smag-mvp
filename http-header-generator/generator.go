package http_header_generator

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type HttpHeaderGenerator struct {
	browserAgent BrowserAgent
}

func New() *HttpHeaderGenerator {
	httpHeader := &HttpHeaderGenerator{}
	data, err := ioutil.ReadFile("useragents.json")
	if err != nil {
		panic(err)
	}
	var userAgent BrowserAgent
	errJson := json.Unmarshal(data, &userAgent)

	if errJson != nil {
		panic(errJson)
	}
	httpHeader.browserAgent = userAgent
	return httpHeader
}

type BrowserAgent []struct {
	UserAgents string `json:"useragent"`
}

func (h *HttpHeaderGenerator) AddHeaders(header *http.Header) {
	header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	header.Add("Accept-Charset", "utf-8")
	header.Add("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	header.Add("Cache-Control", "no-cache")
	header.Add("Content-Type", "application/json; charset=utf-8")
	header.Add("User-Agent", h.GetRandomUserAgent())
}

func (h *HttpHeaderGenerator) GetRandomUserAgent() string {
	randomNumber := rand.Intn(len(h.browserAgent))
	return h.browserAgent[randomNumber].UserAgents
}
