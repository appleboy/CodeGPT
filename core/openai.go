package core

import (
	"context"
	"strconv"

	"github.com/sashabaranov/go-openai"
)

// Usage represents the token usage statistics for a given request.
// It includes counts for the prompt tokens, the completion tokens, and the overall total tokens.
// In addition, it may provide detailed breakdowns for both prompt and completion tokens,
// allowing for deeper insights into token distribution if the corresponding details are available.
type Usage struct {
	PromptTokens            int
	CompletionTokens        int
	TotalTokens             int
	PromptTokensDetails     *openai.PromptTokensDetails
	CompletionTokensDetails *openai.CompletionTokensDetails
}

func (u Usage) String() string {
	s := "Prompt tokens: " + strconv.Itoa(u.PromptTokens)
	if u.PromptTokensDetails != nil && u.PromptTokensDetails.CachedTokens > 0 {
		s += " (CachedTokens: " + strconv.Itoa(u.PromptTokensDetails.CachedTokens) + ")"
	}
	s += ", Completion tokens: " + strconv.Itoa(u.CompletionTokens)
	if u.CompletionTokensDetails != nil && u.CompletionTokensDetails.ReasoningTokens > 0 {
		s += " (ReasoningTokens: " + strconv.Itoa(u.CompletionTokensDetails.ReasoningTokens) + ")"
	}
	s += ", Total tokens: " + strconv.Itoa(u.TotalTokens)
	return s
}

// Response represents the structure of a response from the OpenAI API.
// It contains the content of the response and the usage information.
type Response struct {
	Content string
	Usage   Usage
}

// Generative defines an interface for generative AI operations.
// It includes methods for creating completions and obtaining summary prefixes.
type Generative interface {
	// Completion generates a completion based on the provided content.
	// It takes a context and a string as input and returns a Response pointer and an error.
	Completion(ctx context.Context, content string) (resp *Response, err error)

	// GetSummaryPrefix generates a summary prefix based on the provided content.
	// It takes a context and a string as input and returns a Response pointer and an error.
	GetSummaryPrefix(ctx context.Context, content string) (resp *Response, err error)
}
