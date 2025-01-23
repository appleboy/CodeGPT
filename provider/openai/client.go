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
// It adds the headers from DefaultHeaderTransport to the request before sending it.
// Usage:
//   transport := &DefaultHeaderTransport{
//       Origin: http.DefaultTransport,
//       Header: http.Header{
//           "Authorization": {"Bearer token"},
//       },
//   }
//   client := &http.Client{Transport: transport}
//   resp, err := client.Get("https://example.com")
func (t *DefaultHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, values := range t.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return t.Origin.RoundTrip(req)
}

// NewHeaders creates a new http.Header from the given slice of headers.
// Each header in the slice should be in the format "key=value".
// If a header is not in the correct format, it is skipped.
// Usage:
//   headers := []string{"Authorization=Bearer token", "Content-Type=application/json"}
//   httpHeaders := NewHeaders(headers)
//   fmt.Println(httpHeaders.Get("Authorization")) // Output: Bearer token
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
