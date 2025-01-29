package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

// convertHeaders converts a slice of strings representing HTTP headers
// into an http.Header map. Each string in the input slice should be in
// the format "key=value". If a string cannot be split into exactly two
// parts, or if either the key or value is empty after trimming whitespace,
// that string is ignored.
//
// Parameters:
//
//	headers []string - A slice of strings where each string represents an
//	                   HTTP header in the format "key=value".
//
// Returns:
//
//	http.Header - A map of HTTP headers where the keys are header names
//	              and the values are header values.
func convertHeaders(headers []string) http.Header {
	h := make(http.Header)
	for _, header := range headers {
		// split header into key and value with = as delimiter
		parts := strings.SplitN(header, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			continue
		}
		h.Add(key, value)
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
