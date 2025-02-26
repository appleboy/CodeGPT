package core

import (
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestUsageString(t *testing.T) {
	tests := []struct {
		name     string
		usage    Usage
		expected string
	}{
		{
			name: "without details",
			usage: Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
			expected: "Prompt tokens: 10, Completion tokens: 20, Total tokens: 30",
		},
		{
			name: "with cached tokens",
			usage: Usage{
				PromptTokens:        10,
				CompletionTokens:    20,
				TotalTokens:         30,
				PromptTokensDetails: &openai.PromptTokensDetails{CachedTokens: 5},
			},
			expected: "Prompt tokens: 10 (CachedTokens: 5), Completion tokens: 20, Total tokens: 30",
		},
		{
			name: "with reasoning tokens",
			usage: Usage{
				PromptTokens:            10,
				CompletionTokens:        20,
				TotalTokens:             30,
				CompletionTokensDetails: &openai.CompletionTokensDetails{ReasoningTokens: 7},
			},
			expected: "Prompt tokens: 10, Completion tokens: 20 (ReasoningTokens: 7), Total tokens: 30",
		},
		{
			name: "with both details",
			usage: Usage{
				PromptTokens:            15,
				CompletionTokens:        25,
				TotalTokens:             40,
				PromptTokensDetails:     &openai.PromptTokensDetails{CachedTokens: 3},
				CompletionTokensDetails: &openai.CompletionTokensDetails{ReasoningTokens: 9},
			},
			expected: "Prompt tokens: 15 (CachedTokens: 3), Completion tokens: 25 (ReasoningTokens: 9), Total tokens: 40",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.usage.String()
			if result != tc.expected {
				t.Errorf("Usage.String() = %q, expected %q", result, tc.expected)
			}
		})
	}
}
