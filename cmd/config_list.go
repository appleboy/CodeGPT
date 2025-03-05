package cmd

import (
	"sort"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configCmd.AddCommand(configListCmd)
}

// availableKeys is a list of available config keys
// availableKeys is a map of configuration keys and their descriptions
var availableKeys = map[string]string{
	"git.diff_unified":         "Generate diffs with <n> lines of context, default is 3",
	"git.exclude_list":         "Exclude file from git diff command",
	"git.template_file":        "Template file for commit message",
	"git.template_string":      "Template string for commit message",
	"openai.socks":             "SOCKS proxy",
	"openai.api_key":           "OpenAI API key",
	"openai.model":             "OpenAI model",
	"openai.org_id":            "OpenAI requesting organization",
	"openai.proxy":             "HTTP proxy",
	"output.lang":              "Summarizing language, defaults to English",
	"openai.base_url":          "API base URL to use",
	"openai.timeout":           "Request timeout",
	"openai.max_tokens":        "Maximum number of tokens to generate in the chat completion",
	"openai.temperature":       "Sampling temperature to use, between 0 and 2. Higher values like 0.8 make the output more random, while lower values like 0.2 make it more focused and deterministic",
	"openai.provider":          "Service provider, supports 'openai' or 'azure'",
	"openai.skip_verify":       "Skip verifying TLS certificate",
	"openai.headers":           "Custom headers for OpenAI request",
	"openai.api_version":       "OpenAI API version",
	"openai.top_p":             "Nucleus sampling probability mass. For example, 0.1 means only the tokens comprising the top 10% probability mass are considered",
	"openai.frequency_penalty": "Penalty for new tokens based on their existing frequency in the text so far. Decreases the model's likelihood to repeat the same line verbatim",
	"openai.presence_penalty":  "Penalty for new tokens based on whether they appear in the text so far. Increases the model's likelihood to talk about new topics",
	"prompt.folder":            "Prompt template folder",
}

// configListCmd represents the command to list the configuration values.
// It creates a table with the header "Key" and "Value" and adds the configuration keys and values to the table.
// The api key is hidden for security purposes.
// Finally, it prints the table.
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "show the config list",
	Run: func(cmd *cobra.Command, args []string) {
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create a new table with the header "Key" and "Value"
		tbl := table.New("Key", "Value")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Sort the keys
		keys := make([]string, 0, len(availableKeys))
		for key := range availableKeys {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		// Add the key and value to the table
		for _, v := range keys {
			// Hide the api key
			if v == "openai.api_key" {
				tbl.AddRow(v, "****************")
				continue
			}
			tbl.AddRow(v, viper.Get(v))
		}

		// Print the table
		tbl.Print()
	},
}
