package core

// Platform represents a type for different AI platforms.
type Platform string

const (
	// OpenAI represents the OpenAI platform.
	OpenAI Platform = "openai"
	// Azure represents the Azure platform.
	Azure Platform = "azure"
	// Gemini represents the Gemini platform.
	Gemini Platform = "gemini"
	// Anthropic represents the Anthropic platform.
	Anthropic Platform = "anthropic"
)

// String returns the string representation of the Platform.
func (p Platform) String() string {
	return string(p)
}

// IsValid returns true if the Platform is valid.
func (p Platform) IsValid() bool {
	switch p {
	case OpenAI, Azure, Gemini, Anthropic:
		return true
	}
	return false
}
