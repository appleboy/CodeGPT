package cmd

import (
	"errors"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configCmd.AddCommand(configSetCmd)
	configSetCmd.Flags().StringP("base_url", "b", "", availableKeys["openai.base_url"])
	configSetCmd.Flags().StringP("api_key", "k", "", availableKeys["openai.api_key"])
	configSetCmd.Flags().StringP("model", "m", "gpt-3.5-turbo", availableKeys["openai.model"])
	configSetCmd.Flags().StringP("lang", "l", "en", availableKeys["openai.lang"])
	configSetCmd.Flags().StringP("org_id", "o", "", availableKeys["openai.org_id"])
	configSetCmd.Flags().StringP("proxy", "", "", availableKeys["openai.proxy"])
	configSetCmd.Flags().StringP("socks", "", "", availableKeys["openai.socks"])
	configSetCmd.Flags().DurationP("timeout", "t", defaultTimeout, availableKeys["openai.timeout"])
	configSetCmd.Flags().StringP("template_file", "", "", availableKeys["git.template_file"])
	configSetCmd.Flags().StringP("template_string", "", "", availableKeys["git.template_string"])
	configSetCmd.Flags().IntP("diff_unified", "", 3, availableKeys["git.diff_unified"])
	configSetCmd.Flags().StringP("exclude_list", "", "", availableKeys["git.exclude_list"])
	configSetCmd.Flags().IntP("max_tokens", "", 300, availableKeys["openai.max_tokens"])
	configSetCmd.Flags().Float32P("temperature", "", 1.0, availableKeys["openai.temperature"])
	configSetCmd.Flags().Float32P("top_p", "", 1.0, availableKeys["openai.top_p"])
	configSetCmd.Flags().Float32P("frequency_penalty", "", 0.0, availableKeys["openai.frequency_penalty"])
	configSetCmd.Flags().Float32P("presence_penalty", "", 0.0, availableKeys["openai.presence_penalty"])
	configSetCmd.Flags().StringP("provider", "", "openai", availableKeys["openai.provider"])
	configSetCmd.Flags().StringP("model_name", "", "", availableKeys["openai.model_name"])
	configSetCmd.Flags().BoolP("skip_verify", "", false, availableKeys["openai.skip_verify"])
	configSetCmd.Flags().StringP("headers", "", "", availableKeys["openai.headers"])
	configSetCmd.Flags().StringP("api_version", "", "", availableKeys["openai.api_version"])
	_ = viper.BindPFlag("openai.base_url", configSetCmd.Flags().Lookup("base_url"))
	_ = viper.BindPFlag("openai.org_id", configSetCmd.Flags().Lookup("org_id"))
	_ = viper.BindPFlag("openai.api_key", configSetCmd.Flags().Lookup("api_key"))
	_ = viper.BindPFlag("openai.model", configSetCmd.Flags().Lookup("model"))
	_ = viper.BindPFlag("openai.proxy", configSetCmd.Flags().Lookup("proxy"))
	_ = viper.BindPFlag("openai.socks", configSetCmd.Flags().Lookup("socks"))
	_ = viper.BindPFlag("openai.timeout", configSetCmd.Flags().Lookup("timeout"))
	_ = viper.BindPFlag("openai.max_tokens", configSetCmd.Flags().Lookup("max_tokens"))
	_ = viper.BindPFlag("openai.temperature", configSetCmd.Flags().Lookup("temperature"))
	_ = viper.BindPFlag("output.lang", configSetCmd.Flags().Lookup("lang"))
	_ = viper.BindPFlag("git.diff_unified", configSetCmd.Flags().Lookup("diff_unified"))
	_ = viper.BindPFlag("git.exclude_list", configSetCmd.Flags().Lookup("exclude_list"))
	_ = viper.BindPFlag("git.template_file", configSetCmd.Flags().Lookup("template_file"))
	_ = viper.BindPFlag("git.template_string", configSetCmd.Flags().Lookup("template_string"))
	_ = viper.BindPFlag("openai.provider", configSetCmd.Flags().Lookup("provider"))
	_ = viper.BindPFlag("openai.model_name", configSetCmd.Flags().Lookup("model_name"))
	_ = viper.BindPFlag("openai.skip_verify", configSetCmd.Flags().Lookup("skip_verify"))
	_ = viper.BindPFlag("openai.headers", configSetCmd.Flags().Lookup("headers"))
	_ = viper.BindPFlag("openai.api_version", configSetCmd.Flags().Lookup("api_version"))
}

// configSetCmd updates the config value.
// It takes at least two arguments, the first one being the key and the second one being the value.
// If the key is not available, it returns an error message.
// If the key is "git.exclude_list", it sets the value as a slice of strings.
// It writes the config to file and prints a success message with the config file location.
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "update the config value",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if key is available
		if _, ok := availableKeys[args[0]]; !ok {
			return errors.New("config key is not available, please use `codegpt config list` to see the available keys")
		}

		// Set config value in viper
		if args[0] == "git.exclude_list" {
			viper.Set(args[0], strings.Split(args[1], ","))
		} else {
			viper.Set(args[0], args[1])
		}

		// Write config to file
		if err := viper.WriteConfig(); err != nil {
			return err
		}

		// Print success message with config file location
		color.Green("you can see the config file: %s", viper.ConfigFileUsed())
		return nil
	},
}
