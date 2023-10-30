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
