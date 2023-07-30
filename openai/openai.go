package openai

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/net/proxy"
)

// DefaultModel is the default OpenAI model to use if one is not provided.
var DefaultModel = openai.GPT3Dot5Turbo

// modelMaps maps model names to their corresponding model ID strings.
var modelMaps = map[string]string{
	"gpt-4-32k-0613":         openai.GPT432K0613,
	"gpt-4-32k-0314":         openai.GPT432K0314,
	"gpt-4-32k":              openai.GPT432K,
	"gpt-4-0613":             openai.GPT40613,
	"gpt-4-0314":             openai.GPT40314,
	"gpt-4":                  openai.GPT4,
	"gpt-3.5-turbo-0613":     openai.GPT3Dot5Turbo0613,
	"gpt-3.5-turbo-0301":     openai.GPT3Dot5Turbo0301,
	"gpt-3.5-turbo-16k":      openai.GPT3Dot5Turbo16K,
	"gpt-3.5-turbo-16k-0613": openai.GPT3Dot5Turbo16K0613,
	"gpt-3.5-turbo":          openai.GPT3Dot5Turbo,
	"text-davinci-003":       openai.GPT3TextDavinci003,
	"text-davinci-002":       openai.GPT3TextDavinci002,
	"text-davinci-001":       openai.GPT3TextDavinci001,
	"text-curie-001":         openai.GPT3TextCurie001,
	"text-babbage-001":       openai.GPT3TextBabbage001,
	"text-ada-001":           openai.GPT3TextAda001,
	"davinci-instruct-beta":  openai.GPT3DavinciInstructBeta,
	"davinci":                openai.GPT3Davinci,
	"curie-instruct-beta":    openai.GPT3CurieInstructBeta,
	"curie":                  openai.GPT3Curie,
	"ada":                    openai.GPT3Ada,
	"babbage":                openai.GPT3Babbage,
}

// GetModel returns the model ID corresponding to the given model name.
// If the model name is not recognized, it returns the default model ID.
func GetModel(model string) string {
	v, ok := modelMaps[model]
	if !ok {
		return DefaultModel
	}
	return v
}

// Client is a struct that represents an OpenAI client.
type Client struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
}

type Response struct {
	Content string
	Usage   openai.Usage
}

// CreateChatCompletion is an API call to create a completion for a chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	content string,
) (resp openai.ChatCompletionResponse, err error) {
	req := openai.ChatCompletionRequest{
		Model:       c.model,
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
		TopP:        1,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}

	return c.client.CreateChatCompletion(ctx, req)
}

// CreateCompletion is an API call to create a completion.
// This is the main endpoint of the API. It returns new text, as well as, if requested,
// the probabilities over each alternative token at each position.
//
// If using a fine-tuned model, simply provide the model's ID in the CompletionRequest object,
// and the server will use the model's parameters to generate the completion.
func (c *Client) CreateCompletion(
	ctx context.Context,
	content string,
) (resp openai.CompletionResponse, err error) {
	req := openai.CompletionRequest{
		Model:       c.model,
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
		TopP:        1,
		Prompt:      content,
	}

	return c.client.CreateCompletion(ctx, req)
}

// Completion is a method on the Client struct that takes a context.Context and a string argument
// and returns a string and an error.
func (c *Client) Completion(
	ctx context.Context,
	content string,
) (*Response, error) {
	resp := &Response{}
	switch c.model {
	case openai.GPT3Dot5Turbo,
		openai.GPT3Dot5Turbo0301,
		openai.GPT3Dot5Turbo0613,
		openai.GPT3Dot5Turbo16K,
		openai.GPT3Dot5Turbo16K0613,
		openai.GPT4,
		openai.GPT40314,
		openai.GPT40613,
		openai.GPT432K,
		openai.GPT432K0314,
		openai.GPT432K0613:
		r, err := c.CreateChatCompletion(ctx, content)
		if err != nil {
			return nil, err
		}
		resp.Content = r.Choices[0].Message.Content
		resp.Usage = r.Usage
	default:
		r, err := c.CreateCompletion(ctx, content)
		if err != nil {
			return nil, err
		}
		resp.Content = r.Choices[0].Text
		resp.Usage = r.Usage
	}
	return resp, nil
}

// New creates a new OpenAI API client with the given options.
func New(opts ...Option) (*Client, error) {
	// Create a new config object with the given options.
	cfg := newConfig(opts...)

	// Validate the config object, returning an error if it is invalid.
	if err := cfg.valid(); err != nil {
		return nil, err
	}

	// Create a new client instance with the necessary fields.
	instance := &Client{
		model:       modelMaps[cfg.model],
		maxTokens:   cfg.maxTokens,
		temperature: cfg.temperature,
	}

	// Create a new OpenAI config object with the given API token and other optional fields.
	c := openai.DefaultConfig(cfg.token)
	if cfg.orgID != "" {
		c.OrgID = cfg.orgID
	}
	if cfg.baseURL != "" {
		c.BaseURL = cfg.baseURL
	}

	// Create a new HTTP transport.
	tr := &http.Transport{}
	if cfg.skipVerify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Create a new HTTP client with the specified timeout and proxy, if any.
	httpClient := &http.Client{
		Timeout:   cfg.timeout,
		Transport: tr,
	}

	if cfg.proxyURL != "" {
		proxyURL, _ := url.Parse(cfg.proxyURL)
		tr.Proxy = http.ProxyURL(proxyURL)
	} else if cfg.socksURL != "" {
		dialer, err := proxy.SOCKS5("tcp", cfg.socksURL, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("can't connect to the proxy: %s", err)
		}
		tr.Dial = dialer.Dial
	}

	// Set the OpenAI client to use the default configuration with Azure-specific options, if the provider is Azure.
	if cfg.provider == AZURE {
		defaultAzureConfig := openai.DefaultAzureConfig(cfg.token, cfg.baseURL)
		defaultAzureConfig.AzureModelMapperFunc = func(model string) string {
			return cfg.modelName
		}
		// Set the HTTP client to the one with the specified options.
		defaultAzureConfig.HTTPClient = httpClient
		instance.client = openai.NewClientWithConfig(
			defaultAzureConfig,
		)
	} else {
		// Otherwise, set the OpenAI client to use the HTTP client with the specified options.
		c.HTTPClient = httpClient
		instance.client = openai.NewClientWithConfig(c)
	}

	// Return the resulting client instance.
	return instance, nil
}
