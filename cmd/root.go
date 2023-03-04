package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/appleboy/com/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Short:        "A git prepare-commit-msg hook using ChatGPT",
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
}

// Used for flags.
var (
	cfgFile string
)

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		configFolder := path.Join(home, ".config", "codegpt")
		viper.AddConfigPath(configFolder)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".codegpt")
		cfgFile = path.Join(configFolder, ".codegpt.yaml")

		if !file.IsDir(configFolder) {
			if err := os.MkdirAll(configFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			_, err := os.Create(cfgFile)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Config file was found but another error was produced
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func Execute(ctx context.Context) {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.codegpt.yaml)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)

	// hide completion command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
