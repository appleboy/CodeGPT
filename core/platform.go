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
	// LiteLLM represents the LiteLLM AI gateway platform.
	LiteLLM Platform = "litellm"
	// OpenLLM represents the OpenLLM gateway platform.
	OpenLLM Platform = "openllm"
)

// String returns the string representation of the Platform.
func (p Platform) String() string {
	return string(p)
}

// IsValid returns true if the Platform is valid.
func (p Platform) IsValid() bool {
	switch p {
	case OpenAI, Azure, Gemini, Anthropic, LiteLLM, OpenLLM:
		return true
	}
	return false
}
