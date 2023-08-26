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
		},
		Required: []string{"prefix"},
	},
}

// SummaryPrefixParams is a struct that stores configuration options for the get_summary_prefix function.
type SummaryPrefixParams struct {
	Prefix string `json:"prefix"`
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
