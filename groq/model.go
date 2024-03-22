package groq

type Model string

const (
	LLaMA270bChat          Model = "LLaMA2-70b-chat"
	Mixtral8x7bInstructV01 Model = "Mixtral-8x7b-Instruct-v0.1"
	Gemma7bIt              Model = "Gemma-7b-it"
)

func (m Model) String() string {
	return string(m)
}

func (m Model) GetModel() string {
	return GetModel(m)
}

func (m Model) IsVaild() bool {
	switch m {
	case LLaMA270bChat, Mixtral8x7bInstructV01, Gemma7bIt:
		return true
	default:
		return false
	}
}

var model = map[Model]string{
	LLaMA270bChat:          "llama2-70b-4096",
	Mixtral8x7bInstructV01: "mixtral-8x7b-32768",
	Gemma7bIt:              "gemma-7b-it",
}

// GetModel returns the model name.
func GetModel(modelName Model) string {
	if _, ok := model[modelName]; !ok {
		return model[LLaMA270bChat]
	}
	return model[modelName]
}
