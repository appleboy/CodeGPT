package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the command for custom configuration,
// including settings like openai.api_key, openai.model, and others.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Customize configuration settings (API key, model selection, etc.)",
}
