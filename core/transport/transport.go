package transport

import (
	"net/http"
)

// DefaultHeaderTransport is an http.RoundTripper that adds the given headers to each request,
// and always injects x-app-name and x-app-version headers.
type DefaultHeaderTransport struct {
	Origin     http.RoundTripper
	Header     http.Header
	AppName    string
	AppVersion string
}

// RoundTrip implements the http.RoundTripper interface.
func (t *DefaultHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, values := range t.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	if t.AppName != "" {
		req.Header.Set("x-app-name", t.AppName)
	}
	if t.AppVersion != "" {
		req.Header.Set("x-app-version", t.AppVersion)
	}
	return t.Origin.RoundTrip(req)
}
