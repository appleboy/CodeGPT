package litellm

import (
	"encoding/json"

	"github.com/appleboy/com/bytesconv"
	"github.com/sashabaranov/go-openai/jsonschema"

	openai "github.com/sashabaranov/go-openai"
)

// summaryPrefixFunc defines the OpenAI function-calling schema for extracting
// a conventional commit prefix and scope from a diff summary.
var summaryPrefixFunc = openai.FunctionDefinition{
	Name: "get_summary_prefix",
	Parameters: jsonschema.Definition{
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
			},
			"scope": {
				Type:        jsonschema.String,
				Description: "A short lowercase word identifying the module, package, or component most central to the change",
			},
		},
		Required: []string{"prefix", "scope"},
	},
}

type summaryPrefixParams struct {
	Prefix string `json:"prefix"`
	Scope  string `json:"scope"`
}

func getSummaryPrefixArgs(data string) summaryPrefixParams {
	var prefix summaryPrefixParams
	err := json.Unmarshal(bytesconv.StrToBytes(data), &prefix)
	if err != nil {
		panic(err)
	}
	return prefix
}
