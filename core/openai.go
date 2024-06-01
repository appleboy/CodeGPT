package core

import (
	"context"
)

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
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
