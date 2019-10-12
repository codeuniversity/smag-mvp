package scraper_client

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type SimpleScraperClient struct {
	browserAgent   BrowserAgent
	currentAddress string
	client         *http.Client
	instanceId     string
}

func NewSimpleHttpClient() *SimpleScraperClient {
	client := &SimpleScraperClient{}
	data, err := ioutil.ReadFile("useragents.json")
	if err != nil {
		panic(err)
	}
	var userAgent BrowserAgent
	errJson := json.Unmarshal(data, &userAgent)

	if errJson != nil {
		panic(errJson)
	}
	client.browserAgent = userAgent
	client.client = &http.Client{}
	return client
}

func (s *SimpleScraperClient) WithRetries(times int, f func() error) error {
	var err error
	for i := 0; i < times; i++ {
		err = f()

		if err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return err
}

func (s *SimpleScraperClient) GetClient() *http.Client {
	return s.client
}

func (s *SimpleScraperClient) getRandomUserAgent() string {
	randomNumber := rand.Intn(len(s.browserAgent))
	return s.browserAgent[randomNumber].UserAgents
}

func (s *SimpleScraperClient) AddHeaders(request *http.Request) {
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	request.Header.Add("Accept-Charset", "utf-8")
	request.Header.Add("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("User-Agent", s.getRandomUserAgent())
}
