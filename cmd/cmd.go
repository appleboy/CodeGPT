package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

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
	cfgFile  string
	replacer = strings.NewReplacer("-", "_", ".", "_")
)

const (
	GITHUB = "github"
	DRONE  = "drone"
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/codegpt/.codegpt.yaml)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(reviewCmd)

	// hide completion command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if !file.IsFile(cfgFile) {
			// Config file not found; ignore error if desired
			_, err := os.Create(cfgFile)
			if err != nil {
				log.Fatal(err)
			}
		}
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
	viper.SetEnvKeyReplacer(replacer)

	// Support multiple platforms for CI/CD
	// GitHub Actions need to use `INPUT_` prefix
	// Drone CI need to use `DRONE_` prefix
	switch viper.GetString("platform") {
	case GITHUB:
		viper.SetEnvPrefix("input")
	case DRONE:
		viper.SetEnvPrefix("drone")
	}

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
	if _, err := rootCmd.ExecuteContextC(ctx); err != nil {
		os.Exit(1)
	}
}
