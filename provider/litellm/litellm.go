package litellm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/core/transport"
	"github.com/appleboy/CodeGPT/proxy"
	"github.com/appleboy/CodeGPT/version"

	openai "github.com/sashabaranov/go-openai"
)

var _ core.Generative = (*Client)(nil)

// Client wraps an OpenAI-compatible client configured to talk to a LiteLLM proxy.
// LiteLLM normalizes 100+ provider APIs (Anthropic, Bedrock, Vertex, Groq, etc.)
// into the OpenAI Chat Completions format, so the same go-openai SDK works unchanged.
type Client struct {
	client           *openai.Client
	model            string
	maxTokens        int
	temperature      float32
	topP             float32
	frequencyPenalty float32
	presencePenalty  float32
}

// newBaseRequest builds a ChatCompletionRequest with the client's model parameters.
func (c *Client) newBaseRequest(content string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:               c.model,
		MaxCompletionTokens: c.maxTokens,
		Temperature:         c.temperature,
		TopP:                c.topP,
		FrequencyPenalty:    c.frequencyPenalty,
		PresencePenalty:     c.presencePenalty,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}
}

// convertUsage converts an openai.Usage to a core.Usage.
func convertUsage(u openai.Usage) core.Usage {
	return core.Usage{
		PromptTokens:            u.PromptTokens,
		CompletionTokens:        u.CompletionTokens,
		TotalTokens:             u.TotalTokens,
		CompletionTokensDetails: u.CompletionTokensDetails,
		PromptTokensDetails:     u.PromptTokensDetails,
	}
}

// Completion generates a non-streaming completion via the LiteLLM proxy.
func (c *Client) Completion(ctx context.Context, content string) (*core.Response, error) {
	resp, err := c.client.CreateChatCompletion(ctx, c.newBaseRequest(content))
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from LiteLLM proxy")
	}

	text := resp.Choices[0].Message.Content
	if text == "" && resp.Choices[0].Message.ReasoningContent != "" {
		text = resp.Choices[0].Message.ReasoningContent
	}

	return &core.Response{
		Content: text,
		Usage:   convertUsage(resp.Usage),
	}, nil
}

// CompletionStream streams completion tokens to the writer as they arrive.
func (c *Client) CompletionStream(
	ctx context.Context,
	content string,
	w io.Writer,
) (*core.Response, error) {
	req := c.newBaseRequest(content)
	req.Stream = true
	req.StreamOptions = &openai.StreamOptions{IncludeUsage: true}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	var sb strings.Builder
	var usage openai.Usage
	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		if chunk.Usage != nil {
			usage = *chunk.Usage
		}

		if len(chunk.Choices) > 0 {
			text := chunk.Choices[0].Delta.Content
			if text == "" {
				text = chunk.Choices[0].Delta.ReasoningContent
			}
			if text != "" {
				sb.WriteString(text)
				if _, err := io.WriteString(w, text); err != nil {
					return nil, fmt.Errorf("writing streamed completion: %w", err)
				}
			}
		}
	}

	return &core.Response{
		Content: sb.String(),
		Usage:   convertUsage(usage),
	}, nil
}

// GetSummaryPrefix uses OpenAI-format function calling through the LiteLLM proxy
// to extract a conventional commit prefix. Falls back to plain completion if the
// backing model does not support tool/function calls.
func (c *Client) GetSummaryPrefix(ctx context.Context, content string) (*core.Response, error) {
	req := c.newBaseRequest(content)
	req.Tools = []openai.Tool{
		{
			Type:     openai.ToolTypeFunction,
			Function: &summaryPrefixFunc,
		},
	}
	req.ToolChoice = openai.ToolChoice{
		Type: openai.ToolTypeFunction,
		Function: openai.ToolFunction{
			Name: summaryPrefixFunc.Name,
		},
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		// Function calling not supported by the backing model; fall back to plain completion.
		return c.Completion(ctx, content)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from LiteLLM proxy")
	}

	msg := resp.Choices[0].Message
	usage := convertUsage(resp.Usage)

	if len(msg.ToolCalls) == 0 {
		return &core.Response{
			Content: msg.Content,
			Usage:   usage,
		}, nil
	}

	args := getSummaryPrefixArgs(msg.ToolCalls[len(msg.ToolCalls)-1].Function.Arguments)
	return &core.Response{
		Content: fmt.Sprintf("%s(%s)", args.Prefix, args.Scope),
		Usage:   usage,
	}, nil
}

// New creates a new LiteLLM client that connects to a LiteLLM proxy server.
func New(opts ...Option) (*Client, error) {
	cfg := newConfig(opts...)
	if err := cfg.valid(); err != nil {
		return nil, err
	}

	engine := &Client{
		model:            cfg.model,
		maxTokens:        cfg.maxTokens,
		temperature:      cfg.temperature,
		topP:             cfg.topP,
		frequencyPenalty: cfg.frequencyPenalty,
		presencePenalty:  cfg.presencePenalty,
	}

	c := openai.DefaultConfig(cfg.token)
	c.BaseURL = cfg.baseURL

	httpClient, err := proxy.New(
		proxy.WithProxyURL(cfg.proxyURL),
		proxy.WithSocksURL(cfg.socksURL),
		proxy.WithSkipVerify(cfg.skipVerify),
		proxy.WithTimeout(cfg.timeout),
		proxy.WithHeaders(cfg.headers),
	)
	if err != nil {
		return nil, fmt.Errorf("can't create a new HTTP client: %w", err)
	}

	httpClient.Transport = &transport.DefaultHeaderTransport{
		Origin:     httpClient.Transport,
		Header:     nil,
		AppName:    version.App,
		AppVersion: version.Version,
	}

	c.HTTPClient = httpClient
	engine.client = openai.NewClientWithConfig(c)

	return engine, nil
}
