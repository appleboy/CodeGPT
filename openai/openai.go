package openai

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var DefaultModel = openai.GPT3Dot5Turbo

var modelMaps = map[string]string{
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

func GetModel(model string) string {
	v, ok := modelMaps[model]
	if !ok {
		return DefaultModel
	}
	return v
}

// Clint for OpenAI client interface
type Client struct {
	client *openai.Client
	model  string
}

// CreateChatCompletion — API call to Create a completion for the chat message.
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

// CreateCompletion — API call to create a completion. This is the main endpoint of the API. Returns new text as well
// as, if requested, the probabilities over each alternative token at each position.
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

func (c *Client) Completion(
	ctx context.Context,
	content string,
) (string, error) {
	var message string
	switch c.model {
	case openai.GPT3Dot5Turbo, openai.GPT3Dot5Turbo0301:
		resp, err := c.CreateChatCompletion(ctx, content)
		if err != nil {
			return "", err
		}
		message = resp.Choices[0].Message.Content
	default:
		resp, err := c.CreateCompletion(ctx, content)
		if err != nil {
			return "", err
		}
		message = resp.Choices[0].Text
	}
	return message, nil
}

// New for initialize OpenAI client interface.
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

	if cfg.proxyURL != "" {
		httpClient := &http.Client{
			Timeout: time.Second * 10,
		}
		proxy, _ := url.Parse(cfg.proxyURL)
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		c.HTTPClient = httpClient
	}

	instance.client = openai.NewClientWithConfig(c)

	return instance, nil
}
