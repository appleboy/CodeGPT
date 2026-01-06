package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/provider/anthropic"
	"github.com/appleboy/CodeGPT/provider/gemini"
	"github.com/appleboy/CodeGPT/provider/openai"
	"github.com/appleboy/CodeGPT/util"

	"github.com/spf13/viper"
)

func NewOpenAI(ctx context.Context) (*openai.Client, error) {
	var apiKey string

	// Try to get API key from helper first, fallback to static config
	if helper := viper.GetString("openai.api_key_helper"); helper != "" {
		var refreshInterval time.Duration
		if viper.IsSet("openai.api_key_helper_refresh_interval") {
			// User explicitly set a value (could be 0 to disable cache)
			refreshInterval = time.Duration(
				viper.GetInt("openai.api_key_helper_refresh_interval"),
			) * time.Second
		} else {
			// Not set, use default
			refreshInterval = util.DefaultRefreshInterval
		}
		key, err := util.GetAPIKeyFromHelperWithCache(ctx, helper, refreshInterval)
		if err != nil {
			return nil, err
		}
		apiKey = key
	} else {
		apiKey = viper.GetString("openai.api_key")
	}

	return openai.New(
		openai.WithToken(apiKey),
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
	var apiKey string

	// Try gemini.api_key_helper first
	if helper := viper.GetString("gemini.api_key_helper"); helper != "" {
		var refreshInterval time.Duration
		if viper.IsSet("gemini.api_key_helper_refresh_interval") {
			// User explicitly set a value (could be 0 to disable cache)
			refreshInterval = time.Duration(
				viper.GetInt("gemini.api_key_helper_refresh_interval"),
			) * time.Second
		} else {
			// Not set, use default
			refreshInterval = util.DefaultRefreshInterval
		}
		key, err := util.GetAPIKeyFromHelperWithCache(ctx, helper, refreshInterval)
		if err != nil {
			return nil, err
		}
		apiKey = key
	} else {
		// Fallback to static config: gemini.api_key -> openai.api_key
		apiKey = viper.GetString("gemini.api_key")
		if apiKey == "" {
			// Try openai.api_key_helper as fallback
			if helper := viper.GetString("openai.api_key_helper"); helper != "" {
				var refreshInterval time.Duration
				if viper.IsSet("openai.api_key_helper_refresh_interval") {
					// User explicitly set a value (could be 0 to disable cache)
					refreshInterval = time.Duration(
						viper.GetInt("openai.api_key_helper_refresh_interval"),
					) * time.Second
				} else {
					// Not set, use default
					refreshInterval = util.DefaultRefreshInterval
				}
				key, err := util.GetAPIKeyFromHelperWithCache(ctx, helper, refreshInterval)
				if err != nil {
					return nil, err
				}
				apiKey = key
			} else {
				apiKey = viper.GetString("openai.api_key")
			}
		}
	}

	return gemini.New(
		ctx,
		gemini.WithToken(apiKey),
		gemini.WithModel(viper.GetString("openai.model")),
		gemini.WithMaxTokens(viper.GetInt32("openai.max_tokens")),
		gemini.WithTemperature(float32(viper.GetFloat64("openai.temperature"))),
		gemini.WithTopP(float32(viper.GetFloat64("openai.top_p"))),
		gemini.WithBackend(viper.GetString("gemini.backend")),
		gemini.WithProject(viper.GetString("gemini.project_id")),
		gemini.WithLocation(viper.GetString("gemini.location")),
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
	var apiKey string

	// Try to get API key from helper first, fallback to static config
	if helper := viper.GetString("openai.api_key_helper"); helper != "" {
		var refreshInterval time.Duration
		if viper.IsSet("openai.api_key_helper_refresh_interval") {
			// User explicitly set a value (could be 0 to disable cache)
			refreshInterval = time.Duration(
				viper.GetInt("openai.api_key_helper_refresh_interval"),
			) * time.Second
		} else {
			// Not set, use default
			refreshInterval = util.DefaultRefreshInterval
		}
		key, err := util.GetAPIKeyFromHelperWithCache(ctx, helper, refreshInterval)
		if err != nil {
			return nil, err
		}
		apiKey = key
	} else {
		apiKey = viper.GetString("openai.api_key")
	}

	return anthropic.New(
		anthropic.WithAPIKey(apiKey),
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
		return NewOpenAI(ctx)
	case core.Anthropic:
		return NewAnthropic(ctx)
	}
	return nil, errors.New("invalid provider")
}
