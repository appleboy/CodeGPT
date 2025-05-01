package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/core/transport"
	"github.com/appleboy/CodeGPT/proxy"
	"github.com/appleboy/CodeGPT/version"

	"github.com/appleboy/com/convert"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
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

	usage := core.Usage{
		PromptTokens:     resp.Usage.InputTokens,
		CompletionTokens: resp.Usage.OutputTokens,
		TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
	}

	if resp.Usage.CacheCreationInputTokens > 0 || resp.Usage.CacheReadInputTokens > 0 {
		usage.PromptTokensDetails = &openai.PromptTokensDetails{
			CachedTokens: resp.Usage.CacheCreationInputTokens + resp.Usage.CacheReadInputTokens,
		}
	}

	return &core.Response{
		Content: resp.Content[0].GetText(),
		Usage:   usage,
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

	usage := core.Usage{
		PromptTokens:     resp.Usage.InputTokens,
		CompletionTokens: resp.Usage.OutputTokens,
		TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
	}

	if resp.Usage.CacheCreationInputTokens > 0 || resp.Usage.CacheReadInputTokens > 0 {
		usage.PromptTokensDetails = &openai.PromptTokensDetails{
			CachedTokens: resp.Usage.CacheCreationInputTokens + resp.Usage.CacheReadInputTokens,
		}
	}

	return &core.Response{
		Content: result.Prefix,
		Usage:   usage,
	}, nil
}

// New creates a new Client instance with the provided options.
func New(opts ...Option) (c *Client, err error) {
	// Create a new config object with the given options.
	cfg := newConfig(opts...)

	// Validate the config object, returning an error if it is invalid.
	if err := cfg.valid(); err != nil {
		return nil, err
	}

	httpClient, err := proxy.New(
		proxy.WithProxyURL(cfg.proxyURL),
		proxy.WithSocksURL(cfg.socksURL),
		proxy.WithSkipVerify(cfg.skipVerify),
		proxy.WithTimeout(cfg.timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("can't create a new HTTP client: %w", err)
	}

	// Inject x-app-name and x-app-version headers using core/transport.DefaultHeaderTransport
	httpClient.Transport = &transport.DefaultHeaderTransport{
		Origin:     httpClient.Transport,
		Header:     nil,
		AppName:    version.App,
		AppVersion: version.Version,
	}

	// Create a new client instance with the necessary fields.
	engine := &Client{
		client: anthropic.NewClient(
			cfg.apiKey,
			anthropic.WithHTTPClient(httpClient),
		),
		model:       cfg.model,
		maxTokens:   cfg.maxTokens,
		temperature: cfg.temperature,
		topP:        cfg.topP,
	}

	return engine, nil
}
