package proxy

import (
	"time"
)

// Option is an interface that specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

// optionFunc is a type of function that can be used to implement the Option interface.
// It takes a pointer to a config struct and modifies it.
type optionFunc func(*config)

// Ensure that optionFunc satisfies the Option interface.
var _ Option = (*optionFunc)(nil)

// The apply method of optionFunc type is implemented here to modify the config struct based on the function passed.
func (o optionFunc) apply(c *config) {
	o(c)
}

// WithProxyURL is a function that returns an Option, which sets the proxyURL field of the config struct.
func WithProxyURL(val string) Option {
	return optionFunc(func(c *config) {
		c.proxyURL = val
	})
}

// WithSocksURL is a function that returns an Option, which sets the socksURL field of the config struct.
func WithSocksURL(val string) Option {
	return optionFunc(func(c *config) {
		c.socksURL = val
	})
}

// WithTimeout returns a new Option that sets the timeout for the client configuration.
// It takes a time.Duration value representing the timeout duration.
// It returns an optionFunc that sets the timeout field of the configuration to the provided value.
func WithTimeout(val time.Duration) Option {
	return optionFunc(func(c *config) {
		c.timeout = val
	})
}

// WithHeaders returns a new Option that sets the headers for the http client configuration.
func WithHeaders(headers []string) Option {
	return optionFunc(func(c *config) {
		c.headers = headers
	})
}

// WithSkipVerify returns a new Option that sets the insecure flag for the http client configuration.
func WithSkipVerify(insecure bool) Option {
	return optionFunc(func(c *config) {
		c.insecure = insecure
	})
}

// config is a struct that stores configuration options for the instrumentation.
type config struct {
	proxyURL string
	socksURL string
	timeout  time.Duration
	insecure bool
	headers  []string
}

// newConfig creates a new config object with default values, and applies the given options.
func newConfig(opts ...Option) *config {
	// Create a new config object with default values.
	c := &config{
		timeout:  30 * time.Second,
		insecure: false,
	}

	// Apply each of the given options to the config object.
	for _, opt := range opts {
		opt.apply(c)
	}

	// Return the resulting config object.
	return c
}
