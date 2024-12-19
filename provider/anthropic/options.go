package anthropic

import (
	"errors"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
)

var (
	errorsMissingAPIKey = errors.New("missing api key")
	errorsMissingModel  = errors.New("missing model")
)

var (
	defaultMaxTokens   = 300
	defaultModel       = anthropic.ModelClaude3Haiku20240307
	defaultTemperature = float32(1.0)
	defaultTopP        = float32(1.0)
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

// WithAPIKey is a function that returns an Option, which sets the token field of the config struct.
func WithAPIKey(val string) Option {
	return optionFunc(func(c *config) {
		c.apiKey = val
	})
}

// WithModel is a function that returns an Option, which sets the model field of the config struct.
func WithModel(val string) Option {
	return optionFunc(func(c *config) {
		c.model = anthropic.Model(val)
	})
}

// WithMaxTokens returns a new Option that sets the max tokens for the client configuration.
// The maximum number of tokens to generate in the chat completion.
// The total length of input tokens and generated tokens is limited by the model's context length.
func WithMaxTokens(val int) Option {
	if val <= 0 {
		val = defaultMaxTokens
	}
	return optionFunc(func(c *config) {
		c.maxTokens = val
	})
}

// WithTemperature returns a new Option that sets the temperature for the client configuration.
// What sampling temperature to use, between 0 and 2.
// Higher values like 0.8 will make the output more random,
// while lower values like 0.2 will make it more focused and deterministic.
func WithTemperature(val float32) Option {
	if val <= 0 {
		val = defaultTemperature
	}
	return optionFunc(func(c *config) {
		c.temperature = val
	})
}

// WithTopP returns a new Option that sets the topP for the client configuration.
func WithTopP(val float32) Option {
	return optionFunc(func(c *config) {
		c.topP = val
	})
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

// WithSkipVerify returns a new Option that sets the skipVerify for the client configuration.
func WithSkipVerify(val bool) Option {
	return optionFunc(func(c *config) {
		c.skipVerify = val
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

// config is a struct that stores configuration options for the instrumentation.
type config struct {
	apiKey      string
	model       anthropic.Model
	maxTokens   int
	temperature float32
	topP        float32
	proxyURL    string
	socksURL    string
	skipVerify  bool
	timeout     time.Duration
}

// valid checks whether a config object is valid, returning an error if it is not.
func (cfg *config) valid() error {
	// Check that the token is not empty.
	if cfg.apiKey == "" {
		return errorsMissingAPIKey
	}

	if cfg.model == "" {
		return errorsMissingModel
	}

	// If all checks pass, return nil (no error).
	return nil
}

// newConfig creates a new config object with default values, and applies the given options.
func newConfig(opts ...Option) *config {
	// Create a new config object with default values.
	c := &config{
		model:       defaultModel,
		maxTokens:   defaultMaxTokens,
		temperature: defaultTemperature,
		topP:        defaultTopP,
	}

	// Apply each of the given options to the config object.
	for _, opt := range opts {
		opt.apply(c)
	}

	// Return the resulting config object.
	return c
}
