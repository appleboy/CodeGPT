package litellm

import (
	"reflect"
	"testing"
)

func TestGetSummaryPrefixArgs(t *testing.T) {
	data := `{"prefix": "feat", "scope": "provider", "param2": "value2"}`

	result := getSummaryPrefixArgs(data)

	expected := summaryPrefixParams{
		Prefix: "feat",
		Scope:  "provider",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
