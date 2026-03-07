package gemini

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/core/transport"
	"github.com/appleboy/CodeGPT/version"

	"github.com/appleboy/com/convert"
	"github.com/sashabaranov/go-openai"
	"github.com/yassinebenaid/godump"
	"google.golang.org/genai"
)

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

// CompletionStream streams completion tokens to the writer as they arrive.
func (c *Client) CompletionStream(
	ctx context.Context,
	content string,
	w io.Writer,
) (*core.Response, error) {
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

	var sb strings.Builder
	var writeErr error
	usage := core.Usage{}
	for resp, err := range c.client.Models.GenerateContentStream(ctx, c.model, data, cfg) {
		if err != nil {
			return nil, err
		}

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

		if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
			for _, part := range resp.Candidates[0].Content.Parts {
				if part.Text != "" {
					sb.WriteString(part.Text)
					if writeErr == nil {
						if _, err := io.WriteString(w, part.Text); err != nil {
							writeErr = err
						}
					}
				}
			}
		}
	}

	if writeErr != nil {
		return nil, fmt.Errorf("failed to write streaming response: %w", writeErr)
	}

	return &core.Response{
		Content: sb.String(),
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

	if len(resp.Candidates) == 0 {
		return nil, errors.New("no candidates found")
	}

	cand := resp.Candidates[0]
	if len(cand.Content.Parts) == 0 {
		return nil, errors.New("no content found")
	}

	part := cand.Content.Parts[0]
	if part.FunctionCall == nil || part.FunctionCall.Name != "get_summary_prefix" {
		return nil, errors.New("no function call found")
	}

	prefix, ok := part.FunctionCall.Args["prefix"].(string)
	if !ok || prefix == "" {
		return nil, errors.New("no prefix found")
	}

	scope, _ := part.FunctionCall.Args["scope"].(string)

	if c.debug {
		_ = godump.Dump(resp.Candidates)
	}

	r := &core.Response{
		Content: fmt.Sprintf("%s(%s)", prefix, scope),
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

	var clientConfig *genai.ClientConfig
	switch cfg.backend {
	case genai.BackendVertexAI:
		clientConfig = &genai.ClientConfig{
			HTTPClient: httpClient,
			Backend:    cfg.backend,
			Project:    cfg.projectID,
			Location:   cfg.location,
		}
	case genai.BackendGeminiAPI, genai.BackendUnspecified:
		fallthrough
	default:
		cfg.backend = genai.BackendGeminiAPI
		clientConfig = &genai.ClientConfig{
			APIKey:     cfg.token,
			HTTPClient: httpClient,
			Backend:    cfg.backend,
		}
	}
	client, err := genai.NewClient(ctx, clientConfig)
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
