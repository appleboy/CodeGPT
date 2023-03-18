package git

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

// WithDiffUnified is a function that generate diffs with <n> lines of context instead of the usual three.
func WithDiffUnified(val int) Option {
	return optionFunc(func(c *config) {
		c.diffUnified = val
	})
}

// WithExcludeList returns an Option that sets the excludeList field of a config object to the given value.
func WithExcludeList(val []string) Option {
	return optionFunc(func(c *config) {
		// If the given value is empty, do nothing.
		if len(val) == 0 {
			return
		}
		c.excludeList = val
	})
}

func WithEnableAmend(val bool) Option {
	return optionFunc(func(c *config) {
		c.isAmend = val
	})
}

// config is a struct that stores configuration options for the instrumentation.
type config struct {
	diffUnified int
	excludeList []string
	isAmend     bool
}
