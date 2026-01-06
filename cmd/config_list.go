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

// availableKeys is a map of configuration keys and their descriptions
var availableKeys = map[string]string{
	"git.diff_unified":                       "Number of context lines in git diff output (default: 3)",
	"git.exclude_list":                       "Files to exclude from git diff command",
	"git.template_file":                      "Path to template file for commit messages",
	"git.template_string":                    "Template string for formatting commit messages",
	"openai.socks":                           "SOCKS proxy URL for API connections",
	"openai.api_key":                         "Authentication key for OpenAI API access",
	"openai.api_key_helper":                  "Shell command to dynamically generate API key",
	"openai.api_key_helper_refresh_interval": "Interval in seconds to refresh credentials from apiKeyHelper (default: 900)",
	"openai.model":                           "AI model identifier to use for requests",
	"openai.org_id":                          "Organization ID for multi-org OpenAI accounts",
	"openai.proxy":                           "HTTP proxy URL for API connections",
	"output.lang":                            "Language for summarization output (default: English)",
	"openai.base_url":                        "Custom base URL for API requests",
	"openai.timeout":                         "Maximum duration to wait for API response",
	"openai.max_tokens":                      "Maximum token limit for generated completions",
	"openai.temperature":                     "Randomness control parameter (0-1): lower values for focused results, higher for creative variety",
	"openai.provider":                        "Service provider selection ('openai' or 'azure')",
	"openai.skip_verify":                     "Option to bypass TLS certificate verification",
	"openai.headers":                         "Additional custom HTTP headers for API requests",
	"openai.api_version":                     "Specific API version to target",
	"openai.top_p":                           "Nucleus sampling parameter: controls diversity by limiting to top percentage of probability mass",
	"openai.frequency_penalty":               "Parameter to reduce repetition by penalizing tokens based on their frequency",
	"openai.presence_penalty":                "Parameter to encourage topic diversity by penalizing previously used tokens",
	"prompt.folder":                          "Directory path for custom prompt templates",
	"gemini.project_id":                      "VertexAI project for Gemini provider",
	"gemini.location":                        "VertexAI location for Gemini provider",
	"gemini.backend":                         "Gemini backend (BackendGeminiAPI or BackendVertexAI)",
	"gemini.api_key":                         "API key for Gemini provider",
	"gemini.api_key_helper":                  "Shell command to dynamically generate Gemini API key",
	"gemini.api_key_helper_refresh_interval": "Interval in seconds to refresh Gemini credentials from apiKeyHelper (default: 900)",
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
			if v == "openai.api_key" || v == "gemini.api_key" {
				tbl.AddRow(v, "****************")
				continue
			}
			tbl.AddRow(v, viper.Get(v))
		}

		// Print the table
		tbl.Print()
	},
}
