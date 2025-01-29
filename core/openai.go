package core

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// Usage represents the token usage details for an OpenAI API request.
// It includes the number of tokens used for the prompt, the completion,
// and the total tokens used. Additionally, it may include detailed
// information about the completion tokens.
type Usage struct {
	PromptTokens            int
	CompletionTokens        int
	TotalTokens             int
	CompletionTokensDetails *openai.CompletionTokensDetails
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
