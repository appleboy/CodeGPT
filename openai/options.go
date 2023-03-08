package openai

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

var _ Option = (*optionFunc)(nil)

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func WithToken(val string) Option {
	return optionFunc(func(c *config) {
		c.token = val
	})
}

func WithOrgID(val string) Option {
	return optionFunc(func(c *config) {
		c.orgID = val
	})
}

func WithModel(val string) Option {
	return optionFunc(func(c *config) {
		c.model = val
	})
}

func WithProxyURL(val string) Option {
	return optionFunc(func(c *config) {
		c.proxyURL = val
	})
}

type config struct {
	token    string
	orgID    string
	model    string
	proxyURL string
}
