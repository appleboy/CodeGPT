package anthropic

import (
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// tool represents a structure with a single field Prefix.
// Prefix is a string that is serialized to JSON with the key "prefix".
type tool struct {
	Prefix string `json:"prefix"`
}

// tools is a slice of ToolDefinition that contains the definition for various tools.
// Each ToolDefinition includes the following fields:
// - Name: The name of the tool.
// - Description: A brief description of what the tool does.
// - InputSchema: The schema for the input that the tool expects, defined using jsonschema.Definition.
//   - Type: The type of the input, which is an object.
//   - Properties: A map of property names to their definitions. In this case, it includes:
//   - "prefix": A string that must be one of the specified values ("build", "chore", "ci", "docs", "feat", "fix", "perf", "refactor", "style", "test").
//     This property also has a description indicating that it is the prefix to use for the summary.
//   - Required: A list of required properties, which includes "prefix".
var tools = []anthropic.ToolDefinition{
	{
		Name:        "get_summary_prefix",
		Description: "Get a summary prefix using function call",
		InputSchema: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"prefix": {
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
