package openai

import (
	"context"
	"errors"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

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

func New(token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("missing api key")
	}

	return &Client{
		client: openai.NewClient(token),
	}, nil
}
