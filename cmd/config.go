package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configCmd.PersistentFlags().StringP("api_key", "k", "sk-...", "openai api key")
	configCmd.PersistentFlags().StringP("model", "m", "text-davinci-002", "openai model")
	configCmd.PersistentFlags().StringP("lang", "l", "en", "summarizing language uses English by default")
	viper.BindPFlag("openai.api_key", configCmd.PersistentFlags().Lookup("api_key"))
	viper.BindPFlag("openai.model", configCmd.PersistentFlags().Lookup("model"))
	viper.BindPFlag("output.lang", configCmd.PersistentFlags().Lookup("lang"))
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Add openai config (openai.api_key, openai.model ...)",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "set" {
			fmt.Println("config set key value. ex: config set openai.api_key sk-...")
			os.Exit(1)
		}
		viper.Set(args[1], args[2])
		viper.WriteConfig()
	},
}
