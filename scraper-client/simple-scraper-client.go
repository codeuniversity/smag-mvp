package client

import (
	"net/http"
	"time"

	generator "github.com/codeuniversity/smag-mvp/http_header-generator"
)

// SimpleScraperClient handles retries and setting random headers for scraping
type SimpleScraperClient struct {
	currentAddress string
	client         *http.Client
	instanceID     string
	*generator.HTTPHeaderGenerator
}

// NewSimpleScraperClient returns an initialized SimpleScraperClient
func NewSimpleScraperClient() *SimpleScraperClient {
	client := &SimpleScraperClient{}
	client.HTTPHeaderGenerator = generator.New()
	client.client = &http.Client{}
	return client
}

// WithRetries calls f with retries
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

// Do the request with correct headers
func (s *SimpleScraperClient) Do(request *http.Request) (*http.Response, error) {
	s.AddHeaders(&request.Header)
	return s.client.Do(request)
}
