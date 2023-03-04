package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

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
	version string = "0.0.1"
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
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".codegpt")
		cfgFile = path.Join(home, ".codegpt.yaml")
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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of CodeGPT",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version:", version)
	},
}

func Execute(ctx context.Context) {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.codegpt.yaml)")
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
