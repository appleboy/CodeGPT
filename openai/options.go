package openai

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

// WithOrgID is a function that returns an Option, which sets the orgID field of the config struct.
func WithOrgID(val string) Option {
	return optionFunc(func(c *config) {
		c.orgID = val
	})
}

// WithModel is a function that returns an Option, which sets the model field of the config struct.
func WithModel(val string) Option {
	return optionFunc(func(c *config) {
		c.model = val
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

// WithBaseURL returns a new Option that sets the base URL for the client configuration.
// It takes a string value representing the base URL to use for requests.
// It returns an optionFunc that sets the baseURL field of the configuration to the provided
func WithBaseURL(val string) Option {
	return optionFunc(func(c *config) {
		c.baseURL = val
	})
}

// config is a struct that stores configuration options for the instrumentation.
type config struct {
  baseURL string
	token    string
	orgID    string
	model    string
	proxyURL string
	socksURL string
}
