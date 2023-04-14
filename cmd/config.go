package cmd

import (
	"errors"
	"strings"
	"time"

	"github.com/appleboy/com/array"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var availableKeys = []string{
	"git.diff_unified",
	"git.exclude_list",
	"git.template_file",
	"git.template_string",
	"openai.socks",
	"openai.api_key",
	"openai.model",
	"openai.org_id",
	"openai.proxy",
	"output.lang",
	"openai.base_url",
	"openai.timeout",
	"openai.max_tokens",
	"openai.temperature",

	"openai.provider",
	"openai.model_name",
}

func init() {
	configCmd.PersistentFlags().StringP("base_url", "b", "", "what API base url to use.")
	configCmd.PersistentFlags().StringP("api_key", "k", "", "openai api key")
	configCmd.PersistentFlags().StringP("model", "m", "gpt-3.5-turbo", "openai model")
	configCmd.PersistentFlags().StringP("lang", "l", "en", "summarizing language uses English by default")
	configCmd.PersistentFlags().StringP("org_id", "o", "", "openai requesting organization")
	configCmd.PersistentFlags().StringP("proxy", "", "", "http proxy")
	configCmd.PersistentFlags().StringP("socks", "", "", "socks proxy")
	configCmd.PersistentFlags().DurationP("timeout", "t", 10*time.Second, "http timeout")
	configCmd.PersistentFlags().StringP("template_file", "", "", "git commit message file")
	configCmd.PersistentFlags().StringP("template_string", "", "", "git commit message string")
	configCmd.PersistentFlags().IntP("diff_unified", "", 3, "generate diffs with <n> lines of context, default is 3")
	configCmd.PersistentFlags().IntP("max_tokens", "", 300, "the maximum number of tokens to generate in the chat completion.")
	configCmd.PersistentFlags().Float32P("temperature", "", 0.7, "What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic.")
	configCmd.PersistentFlags().StringSliceP("exclude_list", "", []string{}, "exclude file from `git diff` command")

	configCmd.PersistentFlags().StringP("provider", "", "openai", "srvice provider, only support 'openai' or 'azure'")
	configCmd.PersistentFlags().StringP("model_name", "", "", "model deployment name for Azure cognitive service")

	_ = viper.BindPFlag("openai.base_url", configCmd.PersistentFlags().Lookup("base_url"))
	_ = viper.BindPFlag("openai.org_id", configCmd.PersistentFlags().Lookup("org_id"))
	_ = viper.BindPFlag("openai.api_key", configCmd.PersistentFlags().Lookup("api_key"))
	_ = viper.BindPFlag("openai.model", configCmd.PersistentFlags().Lookup("model"))
	_ = viper.BindPFlag("openai.proxy", configCmd.PersistentFlags().Lookup("proxy"))
	_ = viper.BindPFlag("openai.socks", configCmd.PersistentFlags().Lookup("socks"))
	_ = viper.BindPFlag("openai.timeout", configCmd.PersistentFlags().Lookup("timeout"))
	_ = viper.BindPFlag("openai.max_tokens", configCmd.PersistentFlags().Lookup("max_tokens"))
	_ = viper.BindPFlag("openai.temperature", configCmd.PersistentFlags().Lookup("temperature"))
	_ = viper.BindPFlag("output.lang", configCmd.PersistentFlags().Lookup("lang"))
	_ = viper.BindPFlag("git.diff_unified", configCmd.PersistentFlags().Lookup("diff_unified"))
	_ = viper.BindPFlag("git.exclude_list", configCmd.PersistentFlags().Lookup("exclude_list"))
	_ = viper.BindPFlag("git.template_file", configCmd.PersistentFlags().Lookup("template_file"))
	_ = viper.BindPFlag("git.template_string", configCmd.PersistentFlags().Lookup("template_string"))

	_ = viper.BindPFlag("openai.provider", configCmd.PersistentFlags().Lookup("provider"))
	_ = viper.BindPFlag("openai.model_name", configCmd.PersistentFlags().Lookup("model_name"))
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Add openai config (openai.api_key, openai.model ...)",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] != "set" {
			return errors.New("config set key value. ex: config set openai.api_key sk-...")
		}

		if !array.InSlice(args[1], availableKeys) {
			return errors.New("available key list: " + strings.Join(availableKeys, ", "))
		}

		viper.Set(args[1], args[2])
		if err := viper.WriteConfig(); err != nil {
			return err
		}
		color.Green("you can see the config file: %s", viper.ConfigFileUsed())
		return nil
	},
}
