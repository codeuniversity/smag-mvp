package http_client

import "net/http"

type ClientScraper interface {
	WithRetries(times int, f func() error) error
	GetClient() *http.Client
	AddHeaders(request *http.Request)
}

type HttpStatusError struct {
	s string
}

func (e *HttpStatusError) Error() string {
	return e.s
}
