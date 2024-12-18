package gemini

import (
	"errors"
)

var (
	errorsMissingToken = errors.New("missing gemini api key")
	errorsMissingModel = errors.New("missing model")
)

const (
	defaultMaxTokens   = 300
	defaultModel       = "gemini-1.5-flash-latest"
	defaultTemperature = 1.0
	defaultTopP        = 1.0
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

// WithToken is a function that returns an Option, which sets the token field of the config struct.
func WithToken(val string) Option {
	return optionFunc(func(c *config) {
		c.token = val
	})
}

// WithModel is a function that returns an Option, which sets the model field of the config struct.
func WithModel(val string) Option {
	return optionFunc(func(c *config) {
		c.model = val
	})
}

// WithMaxTokens returns a new Option that sets the max tokens for the client configuration.
// The maximum number of tokens to generate in the chat completion.
// The total length of input tokens and generated tokens is limited by the model's context length.
func WithMaxTokens(val int32) Option {
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

// config is a struct that stores configuration options for the instrumentation.
type config struct {
	token       string
	model       string
	maxTokens   int32
	temperature float32
	topP        float32
}

// valid checks whether a config object is valid, returning an error if it is not.
func (cfg *config) valid() error {
	// Check that the token is not empty.
	if cfg.token == "" {
		return errorsMissingToken
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
