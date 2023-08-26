package openai

import (
	"reflect"
	"testing"
)

func TestGetSummaryPrefixArgs(t *testing.T) {
	data := `{"prefix": "feat", "param2": "value2"}`

	result := GetSummaryPrefixArgs(data)

	expected := SummaryPrefixParams{
		Prefix: "feat",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
