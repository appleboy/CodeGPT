package anthropic

import (
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai/jsonschema"
)

var tools = []anthropic.ToolDefinition{
	{
		Name:        "get_summary_prefix",
		Description: "Get a summary prefix using function call",
		InputSchema: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"unit": {
					Type: jsonschema.String,
					Enum: []string{
						"build", "chore", "ci",
						"docs", "feat", "fix",
						"perf", "refactor", "style",
						"test",
					},
					Description: "The prefix to use for the summary",
				},
			},
			Required: []string{"prefix"},
		},
	},
}
