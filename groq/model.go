package groq

type Model string

const (
	LLaMA38b    Model = "llama3-8b-8192" //
	LLaMA370b   Model = "llama3-70b-8192"
	Mixtral8x7b Model = "mixtral-8x7b-32768"
	Gemma7b     Model = "gemma-7b-it"
)

func (m Model) String() string {
	return string(m)
}

func (m Model) IsVaild() bool {
	switch m {
	case LLaMA38b, LLaMA370b, Mixtral8x7b, Gemma7b:
		return true
	default:
		return false
	}
}
