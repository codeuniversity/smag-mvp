package scraper_client

import (
	"github.com/codeuniversity/smag-mvp/http-header-generator"
	"net/http"
	"time"
)

type SimpleScraperClient struct {
	currentAddress string
	client         *http.Client
	instanceId     string
	*http_header_generator.HttpHeaderGenerator
}

func NewSimpleScraperClient() *SimpleScraperClient {
	client := &SimpleScraperClient{}
	client.HttpHeaderGenerator = http_header_generator.New()
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

func (s *SimpleScraperClient) Do(request *http.Request) (*http.Response, error) {
	s.AddHeaders(&request.Header)
	return s.client.Do(request)
}
