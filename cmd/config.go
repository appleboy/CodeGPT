package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the command for custom configuration,
// including openai.api_key and openai.model and etc...
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "custom config (openai.api_key, openai.model ...)",
}
