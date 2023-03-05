package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/appleboy/com/array"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var availableKeys = []string{"openai.api_key", "openai.model", "openai.org_id", "output.lang"}

func init() {
	configCmd.PersistentFlags().StringP("api_key", "k", "", "openai api key")
	configCmd.PersistentFlags().StringP("model", "m", "gpt-3.5-turbo", "openai model")
	configCmd.PersistentFlags().StringP("lang", "l", "en", "summarizing language uses English by default")
	configCmd.PersistentFlags().StringP("org_id", "o", "", "openai requesting organization")
	_ = viper.BindPFlag("openai.org_id", configCmd.PersistentFlags().Lookup("org_id"))
	_ = viper.BindPFlag("openai.api_key", configCmd.PersistentFlags().Lookup("api_key"))
	_ = viper.BindPFlag("openai.model", configCmd.PersistentFlags().Lookup("model"))
	_ = viper.BindPFlag("output.lang", configCmd.PersistentFlags().Lookup("lang"))
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Add openai config (openai.api_key, openai.model ...)",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "set" {
			log.Fatal("config set key value. ex: config set openai.api_key sk-...")
		}

		if !array.InSlice(args[1], availableKeys) {
			log.Fatal("available key list:", strings.Join(availableKeys, ", "))
		}

		viper.Set(args[1], args[2])
		if err := viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("you can see the config file:", viper.ConfigFileUsed())
	},
}
