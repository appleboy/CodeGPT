package openai

import (
	"encoding/json"

	"github.com/appleboy/com/bytesconv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// SummaryPrefixFunc is a openai function definition.
var SummaryPrefixFunc = openai.FunctionDefinition{
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

// SummaryPrefixParams is a struct that stores configuration options for the get_summary_prefix function.
type SummaryPrefixParams struct {
	Prefix string `json:"prefix"`
	Scope  string `json:"scope"`
}

// GetSummaryPrefixArgs returns the SummaryPrefixParams struct corresponding to the given JSON data.
func GetSummaryPrefixArgs(data string) SummaryPrefixParams {
	var prefix SummaryPrefixParams
	err := json.Unmarshal(bytesconv.StrToBytes(data), &prefix)
	if err != nil {
		panic(err)
	}
	return prefix
}
