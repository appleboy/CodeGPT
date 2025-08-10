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
	cfgFile      string
	promptFolder string
	replacer     = strings.NewReplacer("-", "_", ".", "_")
)

const (
	GITHUB = "github"
	DRONE  = "drone"
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/codegpt/.codegpt.yaml)")
	rootCmd.PersistentFlags().StringVar(&promptFolder, "prompt_folder", "", "prompt folder (default is $HOME/.config/codegpt/prompt)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(CompletionCmd)
	rootCmd.AddCommand(promptCmd)

	// hide completion command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

// initConfig initializes the configuration for the application.
// It sets up the configuration file, environment variables, and prompt folder.
//
// If a configuration file is specified by the cfgFile variable, it uses that file.
// If the file does not exist, it creates a new one.
// If no configuration file is specified, it searches for a configuration file
// named ".codegpt.yaml" in the user's home directory under ".config/codegpt".
//
// The function also sets up environment variable handling to support multiple
// CI/CD platforms, such as GitHub Actions and Drone CI, by setting the appropriate
// environment variable prefixes.
//
// Additionally, it ensures that the prompt folder is correctly set up. If a prompt
// folder is specified by the promptFolder variable, it verifies that it is a directory
// and creates it if it does not exist. If no prompt folder is specified, it defaults
// to a "prompt" directory under the ".config/codegpt" directory in the user's home.
//
// The function uses the Viper library for configuration management and the Cobra
// library for error handling.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		exists, _ := file.IsFile(cfgFile)
		if !exists {
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

		isDir, err := file.IsDir(configFolder)
		if err != nil {
			log.Fatalf("failed to check if config folder %s is a directory: %v", configFolder, err)
		}
		if !isDir {
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

	switch {
	case promptFolder != "":
		// If a prompt folder is specified by the promptFolder variable,
		// check if it is a file. If it is, log a fatal error.
		isFile, err := file.IsFile(promptFolder)
		if err != nil {
			log.Fatalf("failed to check if prompt folder %s is a file: %v", promptFolder, err)
		}
		if isFile {
			log.Fatalf("prompt folder %s is a file", promptFolder)
		}
		// If the prompt folder does not exist, create it.
		isDir, err := file.IsDir(promptFolder)
		if err != nil {
			log.Fatalf("failed to check if prompt folder %s is a directory: %v", promptFolder, err)
		}
		if !isDir {
			if err := os.MkdirAll(promptFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
		// Set the prompt folder in the configuration.
		viper.Set("prompt.folder", promptFolder)
	case viper.GetString("prompt.folder") != "":
		// If the prompt folder is specified in the configuration,
		// retrieve it and check if it is a file. If it is, log a fatal error.
		promptFolder = viper.GetString("prompt.folder")
		isFile, err := file.IsFile(promptFolder)
		if err != nil {
			log.Fatalf("failed to check if prompt folder %s is a file: %v", promptFolder, err)
		}
		if isFile {
			log.Fatalf("prompt folder %s is a file", promptFolder)
		}
		// If the prompt folder does not exist, create it.
		isDir, err := file.IsDir(promptFolder)
		if err != nil {
			log.Fatalf("failed to check if prompt folder %s is a directory: %v", promptFolder, err)
		}
		if !isDir {
			if err := os.MkdirAll(promptFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
	default:
		// If no prompt folder is specified, use the default prompt folder
		// under the ".config/codegpt" directory in the user's home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		targetFolder := path.Join(home, ".config", "codegpt", "prompt")
		isDir, err := file.IsDir(targetFolder)
		if err != nil {
			log.Fatalf("failed to check if target folder %s is a directory: %v", targetFolder, err)
		}
		if !isDir {
			if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
		// Set the prompt folder in the configuration.
		viper.Set("prompt.folder", targetFolder)
	}
}

func Execute(ctx context.Context) {
	if _, err := rootCmd.ExecuteContextC(ctx); err != nil {
		os.Exit(1)
	}
}
