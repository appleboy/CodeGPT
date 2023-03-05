package openai

import (
	"context"
	"errors"

	openai "github.com/sashabaranov/go-openai"
)

// Clint for OpenAI client interface
type Client struct {
	client *openai.Client
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
func New(token, orgID string) (*Client, error) {
	if token == "" {
		return nil, errors.New("missing api key")
	}

	cfg := openai.DefaultConfig(token)
	if orgID != "" {
		cfg.OrgID = orgID
	}

	return &Client{
		client: openai.NewClientWithConfig(cfg),
	}, nil
}
