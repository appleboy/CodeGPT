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
var availableKeys = map[string]string{
	"git.diff_unified":         "generate diffs with <n> lines of context, default is 3",
	"git.exclude_list":         "exclude file from git diff command",
	"git.template_file":        "template file for commit message",
	"git.template_string":      "template string for commit message",
	"openai.socks":             "socks proxy",
	"openai.api_key":           "openai api key",
	"openai.model":             "openai model",
	"openai.org_id":            "openai requesting organization",
	"openai.proxy":             "http proxy",
	"output.lang":              "summarizing language uses English by default",
	"openai.base_url":          "what API base url to use.",
	"openai.timeout":           "http timeout",
	"openai.max_tokens":        "the maximum number of tokens to generate in the chat completion.",
	"openai.temperature":       "What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic.",
	"openai.provider":          "service provider, only support 'openai' or 'azure'",
	"openai.model_name":        "model deployment name for Azure cognitive service",
	"openai.skip_verify":       "skip verify TLS certificate",
	"openai.headers":           "custom headers for openai request",
	"openai.api_version":       "openai api version",
	"openai.top_p":             "An alternative to sampling with temperature, called nucleus sampling, where the model considers the results of the tokens with top_p probability mass. So 0.1 means only the tokens comprising the top 10% probability mass are considered.",
	"openai.frequency_penalty": "Number between 0.0 and 1.0 that penalizes new tokens based on their existing frequency in the text so far. Decreases the model's likelihood to repeat the same line verbatim.",
	"openai.presence_penalty":  "Number between 0.0 and 1.0 that penalizes new tokens based on whether they appear in the text so far. Increases the model's likelihood to talk about new topics.",
}

// configListCmd represents the config list command
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
