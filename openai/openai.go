package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/net/proxy"
)

// DefaultModel is the default OpenAI model to use if one is not provided.
var DefaultModel = openai.GPT3Dot5Turbo

// modelMaps maps model names to their corresponding model ID strings.
var modelMaps = map[string]string{
	"gpt-4-32k-0314":        openai.GPT432K0314,
	"gpt-4-32k":             openai.GPT432K,
	"gpt-4-0314":            openai.GPT40314,
	"gpt-4":                 openai.GPT4,
	"gpt-3.5-turbo":         openai.GPT3Dot5Turbo,
	"gpt-3.5-turbo-0301":    openai.GPT3Dot5Turbo0301,
	"text-davinci-003":      openai.GPT3TextDavinci003,
	"text-davinci-002":      openai.GPT3TextDavinci002,
	"text-davinci-001":      openai.GPT3TextDavinci001,
	"text-curie-001":        openai.GPT3TextCurie001,
	"text-babbage-001":      openai.GPT3TextBabbage001,
	"text-ada-001":          openai.GPT3TextAda001,
	"davinci-instruct-beta": openai.GPT3DavinciInstructBeta,
	"davinci":               openai.GPT3Davinci,
	"curie-instruct-beta":   openai.GPT3CurieInstructBeta,
	"curie":                 openai.GPT3Curie,
	"ada":                   openai.GPT3Ada,
	"babbage":               openai.GPT3Babbage,
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
	client *openai.Client
	model  string
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
		MaxTokens:   200,
		Temperature: 0.7,
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
		MaxTokens:   200,
		Temperature: 0.7,
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
	case openai.GPT3Dot5Turbo, openai.GPT3Dot5Turbo0301:
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

// New is a function that takes a variadic slice of Option types and
// returns a pointer to a Client and an error.
func New(opts ...Option) (*Client, error) {
	cfg := &config{}

	// Loop through each option
	for _, o := range opts {
		// Call the option giving the instantiated
		o.apply(cfg)
	}

	instance := &Client{}
	if cfg.token == "" {
		return nil, errors.New("missing api key")
	}

	v, ok := modelMaps[cfg.model]
	if !ok {
		return nil, errors.New("missing model")
	}
	instance.model = v

	c := openai.DefaultConfig(cfg.token)
	if cfg.orgID != "" {
		c.OrgID = cfg.orgID
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	if cfg.proxyURL != "" {
		proxy, _ := url.Parse(cfg.proxyURL)
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	} else if cfg.socksURL != "" {
		dialer, err := proxy.SOCKS5("tcp", cfg.socksURL, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("can't connect to the proxy: %s", err)
		}
		httpClient.Transport = &http.Transport{
			Dial: dialer.Dial,
		}
	}

	c.HTTPClient = httpClient
	instance.client = openai.NewClientWithConfig(c)

	return instance, nil
}
