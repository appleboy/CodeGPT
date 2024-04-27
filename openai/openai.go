package openai

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/appleboy/CodeGPT/groq"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/net/proxy"
)

// DefaultModel is the default OpenAI model to use if one is not provided.
var DefaultModel = openai.GPT3Dot5Turbo

// modelMaps maps model names to their corresponding model ID strings.
var modelMaps = map[string]string{
	"gpt-4-32k-0613":          openai.GPT432K0613,
	"gpt-4-32k-0314":          openai.GPT432K0314,
	"gpt-4-32k":               openai.GPT432K,
	"gpt-4-0613":              openai.GPT40613,
	"gpt-4-0314":              openai.GPT40314,
	"gpt-4-turbo":             openai.GPT4Turbo,
	"gpt-4-turbo-2024-04-09":  openai.GPT4Turbo20240409,
	"gpt-4-0125-preview":      openai.GPT4Turbo0125,
	"gpt-4-1106-preview":      openai.GPT4Turbo1106,
	"gpt-4-turbo-preview":     openai.GPT4TurboPreview,
	"gpt-4-vision-preview":    openai.GPT4VisionPreview,
	"gpt-4":                   openai.GPT4,
	"gpt-3.5-turbo-0125":      openai.GPT3Dot5Turbo0125,
	"gpt-3.5-turbo-1106":      openai.GPT3Dot5Turbo1106,
	"gpt-3.5-turbo-0613":      openai.GPT3Dot5Turbo0613,
	"gpt-3.5-turbo-0301":      openai.GPT3Dot5Turbo0301,
	"gpt-3.5-turbo-16k":       openai.GPT3Dot5Turbo16K,
	"gpt-3.5-turbo-16k-0613":  openai.GPT3Dot5Turbo16K0613,
	"gpt-3.5-turbo":           openai.GPT3Dot5Turbo,
	"gpt-3.5-turbo-instruct":  openai.GPT3Dot5TurboInstruct,
	"davinci":                 openai.GPT3Davinci,
	"davinci-002":             openai.GPT3Davinci002,
	"curie":                   openai.GPT3Curie,
	"curie-002":               openai.GPT3Curie002,
	"ada":                     openai.GPT3Ada,
	"ada-002":                 openai.GPT3Ada002,
	"babbage":                 openai.GPT3Babbage,
	"babbage-002":             openai.GPT3Babbage002,
	groq.LLaMA38b.String():    groq.LLaMA38b.String(),
	groq.LLaMA370b.String():   groq.LLaMA370b.String(),
	groq.Mixtral8x7b.String(): groq.Mixtral8x7b.String(),
	groq.Gemma7b.String():     groq.Gemma7b.String(),
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
	isFuncCall  bool

	// An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p probability mass.
	// So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	topP float32
	// Number between -2.0 and 2.0.
	// Positive values penalize new tokens based on whether they appear in the text so far,
	// increasing the model's likelihood to talk about new topics.
	presencePenalty float32
	// Number between -2.0 and 2.0.
	// Positive values penalize new tokens based on their existing frequency in the text so far,
	// decreasing the model's likelihood to repeat the same line verbatim.
	frequencyPenalty float32
}

type Response struct {
	Content string
	Usage   openai.Usage
}

// CreateChatCompletion is an API call to create a function call for a chat message.
func (c *Client) CreateFunctionCall(
	ctx context.Context,
	content string,
	f openai.FunctionDefinition,
) (resp openai.ChatCompletionResponse, err error) {
	t := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}

	req := openai.ChatCompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		FrequencyPenalty: c.frequencyPenalty,
		PresencePenalty:  c.presencePenalty,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
		Tools: []openai.Tool{t},
		ToolChoice: openai.ToolChoice{
			Type: openai.ToolTypeFunction,
			Function: openai.ToolFunction{
				Name: f.Name,
			},
		},
	}

	return c.client.CreateChatCompletion(ctx, req)
}

// CreateChatCompletion is an API call to create a completion for a chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	content string,
) (resp openai.ChatCompletionResponse, err error) {
	req := openai.ChatCompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		FrequencyPenalty: c.frequencyPenalty,
		PresencePenalty:  c.presencePenalty,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant.",
			},
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
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		FrequencyPenalty: c.frequencyPenalty,
		PresencePenalty:  c.presencePenalty,
		Prompt:           content,
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
		openai.GPT3Dot5Turbo1106,
		openai.GPT3Dot5Turbo0125,
		openai.GPT4,
		openai.GPT40314,
		openai.GPT40613,
		openai.GPT432K,
		openai.GPT432K0314,
		openai.GPT432K0613,
		openai.GPT4Turbo1106,
		openai.GPT4Turbo0125,
		openai.GPT4TurboPreview,
		openai.GPT4VisionPreview,
		openai.GPT4Turbo,
		openai.GPT4Turbo20240409,
		groq.LLaMA38b.String(),
		groq.LLaMA370b.String(),
		groq.Mixtral8x7b.String(),
		groq.Gemma7b.String():
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
	engine := &Client{
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
		Timeout: cfg.timeout,
	}

	if cfg.proxyURL != "" {
		proxyURL, _ := url.Parse(cfg.proxyURL)
		tr.Proxy = http.ProxyURL(proxyURL)
	} else if cfg.socksURL != "" {
		dialer, err := proxy.SOCKS5("tcp", cfg.socksURL, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("can't connect to the proxy: %s", err)
		}
		tr.DialContext = dialer.(proxy.ContextDialer).DialContext
	}

	// Set the HTTP client to use the default header transport with the specified headers.
	httpClient.Transport = &DefaultHeaderTransport{
		Origin: tr,
		Header: NewHeaders(cfg.headers),
	}

	// Set the OpenAI client to use the default configuration with Azure-specific options, if the provider is Azure.
	if cfg.provider == AZURE {
		defaultAzureConfig := openai.DefaultAzureConfig(cfg.token, cfg.baseURL)
		defaultAzureConfig.AzureModelMapperFunc = func(model string) string {
			return cfg.modelName
		}
		// Set the API version to the one with the specified options.
		if cfg.apiVersion != "" {
			defaultAzureConfig.APIVersion = cfg.apiVersion
		}
		// Set the HTTP client to the one with the specified options.
		defaultAzureConfig.HTTPClient = httpClient
		engine.client = openai.NewClientWithConfig(
			defaultAzureConfig,
		)
	} else {
		// Otherwise, set the OpenAI client to use the HTTP client with the specified options.
		c.HTTPClient = httpClient
		if cfg.apiVersion != "" {
			c.APIVersion = cfg.apiVersion
		}
		engine.client = openai.NewClientWithConfig(c)
	}

	engine.isFuncCall = engine.allowFuncCall(cfg)

	// Return the resulting client engine.
	return engine, nil
}

// allowFuncCall returns true if the model supports function calls.
// https://learn.microsoft.com/en-us/azure/ai-services/openai/how-to/function-calling
// https://platform.openai.com/docs/guides/function-calling/supported-models
// Not all model versions are trained with function calling data.
// Function calling is supported with the following models:
// gpt-4, gpt-4-turbo-preview, gpt-4-0125-preview, gpt-4-1106-preview, gpt-4-0613,
// gpt-3.5-turbo, gpt-3.5-turbo-0125, gpt-3.5-turbo-1106, and gpt-3.5-turbo-0613
// In addition, parallel function calls is supported on the following models:
// gpt-4-turbo-preview, gpt-4-0125-preview, gpt-4-1106-preview,
// gpt-3.5-turbo-0125, and gpt-3.5-turbo-1106
func (c *Client) allowFuncCall(cfg *config) bool {
	if cfg.provider == AZURE && cfg.apiVersion == "2023-07-01-preview" {
		return true
	}

	switch c.model {
	case openai.GPT4Turbo,
		openai.GPT4Turbo20240409,
		openai.GPT4TurboPreview,
		openai.GPT4Turbo0125,
		openai.GPT4Turbo1106,
		openai.GPT40613,
		openai.GPT3Dot5Turbo,
		openai.GPT3Dot5Turbo0125,
		openai.GPT3Dot5Turbo0613,
		openai.GPT3Dot5Turbo1106,
		groq.LLaMA38b.String():
		return true
	default:
		return false
	}
}

// AllowFuncCall returns true if the model supports function calls.
// In an API call, you can describe functions to gpt-3.5-turbo-0613 and gpt-4-0613
// https://platform.openai.com/docs/guides/gpt/chat-completions-api
func (c *Client) AllowFuncCall() bool {
	return c.isFuncCall
}
