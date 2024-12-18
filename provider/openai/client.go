package openai

import (
	"net/http"
	"strings"
)

// DefaultHeaderTransport is an http.RoundTripper that adds the given headers to
type DefaultHeaderTransport struct {
	Origin http.RoundTripper
	Header http.Header
}

// RoundTrip implements the http.RoundTripper interface.
func (t *DefaultHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, values := range t.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return t.Origin.RoundTrip(req)
}

// NewHeaders creates a new http.Header from the given slice of headers.
func NewHeaders(headers []string) http.Header {
	h := make(http.Header)
	for _, header := range headers {
		// split header into key and value with = as delimiter
		vals := strings.Split(header, "=")
		if len(vals) != 2 {
			continue
		}
		h.Add(vals[0], vals[1])
	}
	return h
}
