package transport

import (
	"net/http"
	"strings"
)

// NewHeaders creates a new http.Header from the given slice of headers.
// Each header in the slice should be in the format "key=value".
// If a header is not in the correct format, it is skipped.
func NewHeaders(headers []string) http.Header {
	h := make(http.Header)
	for _, header := range headers {
		vals := strings.Split(header, "=")
		if len(vals) != 2 {
			continue
		}
		h.Add(vals[0], vals[1])
	}
	return h
}
