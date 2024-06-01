package core

type Platform string

const (
	OpenAI Platform = "openai"
	Azure  Platform = "azure"
	Gemini Platform = "gemini"
)

// String returns the string representation of the Platform.
func (p Platform) String() string {
	return string(p)
}

// IsValid returns true if the Platform is valid.
func (p Platform) IsValid() bool {
	switch p {
	case OpenAI, Azure, Gemini:
		return true
	}
	return false
}
