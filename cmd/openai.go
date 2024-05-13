package cmd

import (
	"github.com/appleboy/CodeGPT/openai"

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
		openai.WithProvider(viper.GetString("openai.provider")),
		openai.WithSkipVerify(viper.GetBool("openai.skip_verify")),
		openai.WithHeaders(viper.GetStringSlice("openai.headers")),
		openai.WithAPIVersion(viper.GetString("openai.api_version")),
		openai.WithTopP(float32(viper.GetFloat64("openai.top_p"))),
		openai.WithFrequencyPenalty(float32(viper.GetFloat64("openai.frequency_penalty"))),
		openai.WithPresencePenalty(float32(viper.GetFloat64("openai.presence_penalty"))),
	)
}
