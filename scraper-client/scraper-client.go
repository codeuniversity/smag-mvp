package client

import "net/http"

// ScraperClient is some implementation of a http client usable for scraping
type ScraperClient interface {
	WithRetries(times int, f func() error) error
	Do(request *http.Request) (*http.Response, error)
}

// HTTPStatusError ...
type HTTPStatusError struct {
	S string
}

// Error for the error interface
func (e *HTTPStatusError) Error() string {
	return e.S
}
