package gemini

import "github.com/google/generative-ai-go/genai"

var summaryPrefixFunc = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "get_summary_prefix",
		Description: "Get a summary prefix using function call",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"prefix": {
					Type:        genai.TypeString,
					Description: "The prefix to use for the summary",
					Enum: []string{
						"build", "chore", "ci",
						"docs", "feat", "fix",
						"perf", "refactor", "style",
						"test",
					},
				},
			},
			Required: []string{"prefix"},
		},
	}},
}
