package openai

import (
	"context"
	"errors"

	openai "github.com/sashabaranov/go-openai"
)

var modelMaps = map[string]string{
	"gpt-3.5-turbo":      openai.GPT3Dot5Turbo,
	"gpt-3.5-turbo-0301": openai.GPT3Dot5Turbo0301,
	"text-davinci-003":   openai.GPT3TextDavinci003,
	"text-davinci-002":   openai.GPT3TextDavinci002,
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
	return c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   200,
			Temperature: 0.7,
			TopP:        1,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
}

// New for initialize OpenAI client interface.
func New(token, model, orgID string) (*Client, error) {
	instance := &Client{}
	if token == "" {
		return nil, errors.New("missing api key")
	}

	v, ok := modelMaps[model]
	if !ok {
		return nil, errors.New("missing model")
	}
	instance.model = v

	cfg := openai.DefaultConfig(token)
	if orgID != "" {
		cfg.OrgID = orgID
	}
	instance.client = openai.NewClientWithConfig(cfg)

	return instance, nil
}
