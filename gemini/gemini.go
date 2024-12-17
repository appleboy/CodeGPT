package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/util"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	client      *genai.GenerativeModel
	model       string
	maxTokens   int32
	temperature float32
	topP        float32
	debug       bool
}

// Completion is a method on the Client struct that takes a context.Context and a string argument
func (c *Client) Completion(ctx context.Context, content string) (*core.Response, error) {
	resp, err := c.client.GenerateContent(ctx, genai.Text(content))
	if err != nil {
		return nil, err
	}

	var ret string

	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			ret += fmt.Sprintf("%v", part)
		}
	}

	return &core.Response{
		Content: ret,
		Usage: core.Usage{
			PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
			CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
		},
	}, nil
}

// GetSummaryPrefix is an API call to get a summary prefix using function call.
func (c *Client) GetSummaryPrefix(ctx context.Context, content string) (*core.Response, error) {
	c.client.Tools = []*genai.Tool{summaryPrefixFunc}

	// Start new chat session.
	session := c.client.StartChat()

	// Send the message to the generative model.
	resp, err := session.SendMessage(ctx, genai.Text(content))
	if err != nil {
		return nil, err
	}

	part := resp.Candidates[0].Content.Parts[0]

	r := &core.Response{
		Content: strings.TrimSpace(strings.TrimSuffix(fmt.Sprintf("%v", part), "\n")),
		Usage: core.Usage{
			PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
			CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
		},
	}

	if c.debug {
		// Check that you got the expected function call back.
		funcall, ok := part.(genai.FunctionCall)
		if !ok {
			return nil, fmt.Errorf("expected type FunctionCall, got %T", part)
		}
		if g, e := funcall.Name, summaryPrefixFunc.FunctionDeclarations[0].Name; g != e {
			return nil, fmt.Errorf("expected FunctionCall.Name %q, got %q", e, g)
		}
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

	// Create a new client instance with the necessary fields.
	engine := &Client{
		model:       cfg.model,
		maxTokens:   cfg.maxTokens,
		temperature: cfg.temperature,
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.token))
	if err != nil {
		return nil, err
	}

	engine.client = client.GenerativeModel(engine.model)
	engine.client.MaxOutputTokens = util.Int32Ptr(engine.maxTokens)
	engine.client.Temperature = util.Float32Ptr(engine.temperature)
	engine.client.TopP = util.Float32Ptr(engine.topP)

	return engine, nil
}
