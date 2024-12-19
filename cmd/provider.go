package cmd

import (
	"context"
	"errors"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/provider/anthropic"
	"github.com/appleboy/CodeGPT/provider/gemini"
	"github.com/appleboy/CodeGPT/provider/openai"

	"github.com/spf13/viper"
)

func NewOpenAI() (*openai.Client, error) {
	return openai.New(
		openai.WithToken(viper.GetString("openai.api_key")),
		openai.WithModel(viper.GetString("openai.model")),
		openai.WithOrgID(viper.GetString("openai.org_id")),
		openai.WithProxyURL(viper.GetString("openai.proxy")),
		openai.WithSocksURL(viper.GetString("openai.socks")),
		openai.WithBaseURL(viper.GetString("openai.base_url")),
		openai.WithTimeout(viper.GetDuration("openai.timeout")),
		openai.WithMaxTokens(viper.GetInt("openai.max_tokens")),
		openai.WithTemperature(float32(viper.GetFloat64("openai.temperature"))),
		openai.WithProvider(core.Platform(viper.GetString("openai.provider"))),
		openai.WithSkipVerify(viper.GetBool("openai.skip_verify")),
		openai.WithHeaders(viper.GetStringSlice("openai.headers")),
		openai.WithAPIVersion(viper.GetString("openai.api_version")),
		openai.WithTopP(float32(viper.GetFloat64("openai.top_p"))),
		openai.WithFrequencyPenalty(float32(viper.GetFloat64("openai.frequency_penalty"))),
		openai.WithPresencePenalty(float32(viper.GetFloat64("openai.presence_penalty"))),
	)
}

// NewGemini returns a new Gemini client
func NewGemini(ctx context.Context) (*gemini.Client, error) {
	return gemini.New(
		ctx,
		gemini.WithToken(viper.GetString("openai.api_key")),
		gemini.WithModel(viper.GetString("openai.model")),
		gemini.WithMaxTokens(viper.GetInt32("openai.max_tokens")),
		gemini.WithTemperature(float32(viper.GetFloat64("openai.temperature"))),
		gemini.WithTopP(float32(viper.GetFloat64("openai.top_p"))),
	)
}

// NewAnthropic creates a new instance of the anthropic.Client using configuration
// values retrieved from Viper. The configuration values include the API key,
// model, maximum tokens, temperature, and top_p.
//
// Parameters:
//   - ctx: The context for the client.
//
// Returns:
//   - A pointer to an anthropic.Client instance.
//   - An error if the client could not be created.
func NewAnthropic(ctx context.Context) (*anthropic.Client, error) {
	return anthropic.New(
		anthropic.WithAPIKey(viper.GetString("openai.api_key")),
		anthropic.WithModel(viper.GetString("openai.model")),
		anthropic.WithMaxTokens(viper.GetInt("openai.max_tokens")),
		anthropic.WithTemperature(float32(viper.GetFloat64("openai.temperature"))),
		anthropic.WithTopP(float32(viper.GetFloat64("openai.top_p"))),
		anthropic.WithProxyURL(viper.GetString("openai.proxy")),
		anthropic.WithSocksURL(viper.GetString("openai.socks")),
		anthropic.WithSkipVerify(viper.GetBool("openai.skip_verify")),
		anthropic.WithTimeout(viper.GetDuration("openai.timeout")),
	)
}

// GetClient returns the generative client based on the platform
func GetClient(ctx context.Context, p core.Platform) (core.Generative, error) {
	switch p {
	case core.Gemini:
		return NewGemini(ctx)
	case core.OpenAI, core.Azure:
		return NewOpenAI()
	case core.Anthropic:
		return NewAnthropic(ctx)
	}
	return nil, errors.New("invalid provider")
}
