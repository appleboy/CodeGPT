package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

// convertHeaders takes a slice of strings representing HTTP headers in the format "key=value"
// and converts them into an http.Header type. If a header string does not contain exactly one "=",
// it is ignored. The resulting http.Header map is returned.
func convertHeaders(headers []string) http.Header {
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

// defaultHeaderTransport is a custom implementation of http.RoundTripper
// that allows setting default headers for each request. It wraps an existing
// http.RoundTripper (origin) and adds the specified headers (header) to each
// outgoing request.
type defaultHeaderTransport struct {
	origin http.RoundTripper
	header http.Header
}

// RoundTrip executes a single HTTP transaction and returns
// a Response for the provided Request. It adds custom headers
// from the defaultHeaderTransport to the request before
// delegating the actual round-trip to the original transport.
func (t *defaultHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.origin == nil {
		return nil, fmt.Errorf("origin RoundTripper is nil")
	}
	for key, values := range t.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return t.origin.RoundTrip(req)
}

// New creates a new HTTP client with the provided options.
// It configures the client with optional TLS settings, proxy settings, and custom headers.
//
// Parameters:
//
//	opts - A variadic list of Option functions to configure the client.
//
// Returns:
//
//	*http.Client - A pointer to the configured HTTP client.
//	error - An error if the proxy URL is invalid or if there is an issue connecting to the SOCKS5 proxy.
func New(opts ...Option) (*http.Client, error) {
	cfg := newConfig(opts...)
	if cfg == nil {
		return nil, fmt.Errorf("configuration is nil")
	}

	// Create a new HTTP transport with optional TLS configuration.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.insecure}, //nolint:gosec
	}

	// Create a new HTTP client with the specified timeout.
	httpClient := &http.Client{
		Timeout:   cfg.timeout,
		Transport: tr,
	}

	// Configure proxy settings if provided.
	if cfg.proxyURL != "" {
		proxyURL, err := url.Parse(cfg.proxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %s", err)
		}
		tr.Proxy = http.ProxyURL(proxyURL)
	} else if cfg.socksURL != "" {
		dialer, err := proxy.SOCKS5("tcp", cfg.socksURL, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("can't connect to the SOCKS5 proxy: %s", err)
		}
		tr.DialContext = dialer.(proxy.ContextDialer).DialContext
	}

	// Set the HTTP client to use the default header transport with the specified headers.
	httpClient.Transport = &defaultHeaderTransport{
		origin: tr,
		header: convertHeaders(cfg.headers),
	}

	return httpClient, nil
}
