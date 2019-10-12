package scraper_client

import "net/http"

type ScraperClient interface {
	WithRetries(times int, f func() error) error
	GetClient() *http.Client
	AddHeaders(request *http.Request)
}

type HttpStatusError struct {
	S string
}

func (e *HttpStatusError) Error() string {
	return e.S
}

type BrowserAgent []struct {
	UserAgents string `json:"useragent"`
}
