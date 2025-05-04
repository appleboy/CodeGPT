package gemini

import (
	"context"
	"errors"
	"net/http"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/core/transport"
	"github.com/appleboy/CodeGPT/version"

	"github.com/appleboy/com/convert"
	"github.com/sashabaranov/go-openai"
	"github.com/yassinebenaid/godump"
	"google.golang.org/genai"
)

var ErrInvalidFunctionCall = errors.New("invalid function call")

type Client struct {
	client      *genai.Client
	model       string
	maxTokens   int32
	temperature float32
	topP        float32
	debug       bool
}

// Completion is a method on the Client struct that takes a context.Context and a string argument
func (c *Client) Completion(ctx context.Context, content string) (*core.Response, error) {
	cfg := &genai.GenerateContentConfig{
		TopP:            convert.ToPtr(c.topP),
		Temperature:     convert.ToPtr(c.temperature),
		MaxOutputTokens: c.maxTokens,
	}
	data := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: content,
				},
			},
		},
	}

	resp, err := c.client.Models.GenerateContent(ctx, c.model, data, cfg)
	if err != nil {
		return nil, err
	}

	usage := core.Usage{}
	if resp.UsageMetadata != nil {
		usage.PromptTokens = int(resp.UsageMetadata.PromptTokenCount)
		usage.CompletionTokens = int(resp.UsageMetadata.CandidatesTokenCount)
		usage.TotalTokens = int(resp.UsageMetadata.TotalTokenCount)
		if resp.UsageMetadata.CachedContentTokenCount > 0 {
			usage.PromptTokensDetails = &openai.PromptTokensDetails{
				CachedTokens: int(resp.UsageMetadata.CachedContentTokenCount),
			}
		}
	}

	return &core.Response{
		Content: resp.Text(),
		Usage:   usage,
	}, nil
}

// GetSummaryPrefix is an API call to get a summary prefix using function call.
func (c *Client) GetSummaryPrefix(ctx context.Context, content string) (*core.Response, error) {
	cfg := &genai.GenerateContentConfig{
		MaxOutputTokens: c.maxTokens,
		TopP:            convert.ToPtr(c.topP),
		Temperature:     convert.ToPtr(c.temperature),
		Tools:           []*genai.Tool{summaryPrefixFunc},
		ToolConfig: &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				Mode: genai.FunctionCallingConfigModeAny,
				AllowedFunctionNames: []string{
					"get_summary_prefix",
				},
			},
		},
	}
	data := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: content,
				},
			},
		},
	}

	resp, err := c.client.Models.GenerateContent(ctx, c.model, data, cfg)
	if err != nil {
		return nil, err
	}

	usage := core.Usage{}
	if resp.UsageMetadata != nil {
		usage.PromptTokens = int(resp.UsageMetadata.PromptTokenCount)
		usage.CompletionTokens = int(resp.UsageMetadata.CandidatesTokenCount)
		usage.TotalTokens = int(resp.UsageMetadata.TotalTokenCount)
		if resp.UsageMetadata.CachedContentTokenCount > 0 {
			usage.PromptTokensDetails = &openai.PromptTokensDetails{
				CachedTokens: int(resp.UsageMetadata.CachedContentTokenCount),
			}
		}
	}

	if len(resp.Candidates) == 0 ||
		resp.Candidates[0].Content.Parts[0].FunctionCall.Name != "get_summary_prefix" ||
		resp.Candidates[0].Content.Parts[0].FunctionCall.Args == nil ||
		resp.Candidates[0].Content.Parts[0].FunctionCall.Args["prefix"] == nil ||
		resp.Candidates[0].Content.Parts[0].FunctionCall.Args["prefix"].(string) == "" {
		return nil, ErrInvalidFunctionCall
	}

	if c.debug {
		_ = godump.Dump(resp.Candidates)
	}

	r := &core.Response{
		Content: resp.Candidates[0].Content.Parts[0].FunctionCall.Args["prefix"].(string),
		Usage:   usage,
	}

	return r, nil
}

func New(ctx context.Context, opts ...Option) (c *Client, err error) {
	// Create a new config object with the given options.
	cfg := newConfig(opts...)

	// Validate the config object, returning an error if it is invalid.
	if err := cfg.valid(); err != nil {
		return nil, err
	}

	// Inject x-app-name and x-app-version headers using core/transport.DefaultHeaderTransport
	httpClient := &http.Client{
		Transport: &transport.DefaultHeaderTransport{
			Origin:     http.DefaultTransport,
			Header:     nil,
			AppName:    version.App,
			AppVersion: version.Version,
		},
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     cfg.token,
		HTTPClient: httpClient,
		Backend:    genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	engine := &Client{
		client:      client,
		model:       cfg.model,
		maxTokens:   cfg.maxTokens,
		temperature: cfg.temperature,
		topP:        cfg.topP,
	}

	return engine, nil
}
