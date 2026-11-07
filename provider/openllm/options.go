package openllm

import (
	"errors"
	"time"
)

var errorsMissingModel = errors.New("missing model")

const (
	defaultMaxTokens   = 300
	defaultTemperature = float32(1.0)
	defaultTopP        = float32(1.0)
	defaultBaseURL     = "http://127.0.0.1:8787/v1"
)

// Option is an interface that specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

// optionFunc is a type of function that can be used to implement the Option interface.
type optionFunc func(*config)

var _ Option = (*optionFunc)(nil)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithToken sets the API key for the OpenLLM gateway.
func WithToken(val string) Option {
	return optionFunc(func(c *config) {
		c.token = val
	})
}

// WithModel sets the model identifier. OpenLLM accepts built-in aliases such as
// "ultra", "plus", and "lite", or provider model IDs discovered via /v1/models.
func WithModel(val string) Option {
	return optionFunc(func(c *config) {
		c.model = val
	})
}

// WithBaseURL sets the OpenLLM base URL. Defaults to the local daemon endpoint;
// point it at https://www.openllm.sh/v1 to use the cloud gateway instead.
func WithBaseURL(val string) Option {
	return optionFunc(func(c *config) {
		if val != "" {
			c.baseURL = val
		}
	})
}

// WithProxyURL sets the HTTP proxy URL for outbound connections.
func WithProxyURL(val string) Option {
	return optionFunc(func(c *config) {
		c.proxyURL = val
	})
}

// WithSocksURL sets the SOCKS proxy URL for outbound connections.
func WithSocksURL(val string) Option {
	return optionFunc(func(c *config) {
		c.socksURL = val
	})
}

// WithTimeout sets the per-request HTTP timeout.
func WithTimeout(val time.Duration) Option {
	return optionFunc(func(c *config) {
		c.timeout = val
	})
}

// WithMaxTokens sets the maximum token limit for generated completions.
func WithMaxTokens(val int) Option {
	if val <= 0 {
		val = defaultMaxTokens
	}
	return optionFunc(func(c *config) {
		c.maxTokens = val
	})
}

// WithTemperature sets the sampling temperature (0-2).
func WithTemperature(val float32) Option {
	if val <= 0 {
		val = defaultTemperature
	}
	return optionFunc(func(c *config) {
		c.temperature = val
	})
}

// WithTopP sets the nucleus sampling parameter.
func WithTopP(val float32) Option {
	return optionFunc(func(c *config) {
		c.topP = val
	})
}

// WithFrequencyPenalty sets the frequency penalty parameter.
func WithFrequencyPenalty(val float32) Option {
	return optionFunc(func(c *config) {
		c.frequencyPenalty = val
	})
}

// WithPresencePenalty sets the presence penalty parameter.
func WithPresencePenalty(val float32) Option {
	return optionFunc(func(c *config) {
		c.presencePenalty = val
	})
}

// WithSkipVerify disables TLS certificate verification.
func WithSkipVerify(val bool) Option {
	return optionFunc(func(c *config) {
		c.skipVerify = val
	})
}

// WithHeaders sets additional HTTP headers for requests.
func WithHeaders(headers []string) Option {
	return optionFunc(func(c *config) {
		c.headers = headers
	})
}

type config struct {
	token            string
	model            string
	baseURL          string
	proxyURL         string
	socksURL         string
	timeout          time.Duration
	maxTokens        int
	temperature      float32
	topP             float32
	frequencyPenalty float32
	presencePenalty  float32
	skipVerify       bool
	headers          []string
}

// valid allows an empty token: the local daemon authenticates with the saved
// login when no bearer key is configured.
func (cfg *config) valid() error {
	if cfg.model == "" {
		return errorsMissingModel
	}
	return nil
}

func newConfig(opts ...Option) *config {
	c := &config{
		maxTokens:   defaultMaxTokens,
		temperature: defaultTemperature,
		topP:        defaultTopP,
		baseURL:     defaultBaseURL,
	}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}
