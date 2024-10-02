package core

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type Usage struct {
	PromptTokens            int
	CompletionTokens        int
	TotalTokens             int
	CompletionTokensDetails *openai.CompletionTokensDetails
}

type Response struct {
	Content string
	Usage   Usage
}

type Generative interface {
	// CreateCompletion is an API call to create a completion.
	Completion(ctx context.Context, content string) (resp *Response, err error)
	// GetSummaryPrefix is an API call to get a summary prefix using function call.
	GetSummaryPrefix(ctx context.Context, content string) (resp *Response, err error)
}
