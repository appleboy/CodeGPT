package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/appleboy/CodeGPT/core"

	"github.com/appleboy/com/convert"
	"github.com/liushuangls/go-anthropic/v2"
)

var _ core.Generative = (*Client)(nil)

type Client struct {
	client      *anthropic.Client
	model       anthropic.Model
	maxTokens   int
	temperature float32
	topP        float32
}

// Completion is a method on the Client struct that takes a context.Context and a string argument
func (c *Client) Completion(ctx context.Context, content string) (*core.Response, error) {
	resp, err := c.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: c.model,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(content),
		},
		MaxTokens:   c.maxTokens,
		Temperature: convert.ToPtr(c.temperature),
		TopP:        convert.ToPtr(c.topP),
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
		return nil, err
	}

	return &core.Response{
		Content: resp.Content[0].GetText(),
		Usage: core.Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
	}, nil
}

// GetSummaryPrefix is an API call to get a summary prefix using function call.
func (c *Client) GetSummaryPrefix(ctx context.Context, content string) (*core.Response, error) {
	request := anthropic.MessagesRequest{
		Model: c.model,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(content),
		},
		MaxTokens: c.maxTokens,
		Tools:     tools,
	}

	resp, err := c.client.CreateMessages(ctx, request)
	if err != nil {
		return nil, err
	}

	var toolUse *anthropic.MessageContentToolUse

	for _, c := range resp.Content {
		if c.Type == anthropic.MessagesContentTypeToolUse {
			toolUse = c.MessageContentToolUse
		}
	}

	if toolUse == nil {
		return nil, errors.New("no tool use found in response")
	}

	var result tool
	if err := json.Unmarshal(toolUse.Input, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool use input: %w", err)
	}

	return &core.Response{
		Content: result.Prefix,
		Usage: core.Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
	}, nil
}

func New(opts ...Option) (c *Client, err error) {
	// Create a new config object with the given options.
	cfg := newConfig(opts...)

	// Validate the config object, returning an error if it is invalid.
	if err := cfg.valid(); err != nil {
		return nil, err
	}

	// Create a new client instance with the necessary fields.
	engine := &Client{
		client:      anthropic.NewClient(cfg.apiKey),
		model:       cfg.model,
		maxTokens:   cfg.maxTokens,
		temperature: cfg.temperature,
		topP:        cfg.topP,
	}

	return engine, nil
}
