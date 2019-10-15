package scraper_client

import "net/http"

type ScraperClient interface {
	WithRetries(times int, f func() error) error
	Do(request *http.Request) (*http.Response, error)
}

type HttpStatusError struct {
	S string
}

func (e *HttpStatusError) Error() string {
	return e.S
}
