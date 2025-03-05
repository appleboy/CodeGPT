package cmd

import (
	"os"
	"path"

	"github.com/appleboy/CodeGPT/prompt"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loadPromptData bool

func init() {
	promptCmd.PersistentFlags().BoolVar(&loadPromptData, "load", false,
		"Load default prompt templates into your configuration")
}

var defaultPromptDataKeys = []string{
	prompt.CodeReviewTemplate,
	prompt.SummarizeFileDiffTemplate,
	prompt.SummarizeTitleTemplate,
	prompt.ConventionalCommitTemplate,
}

// promptCmd is a Cobra command to load default prompt data into a specified folder.
// It provides functionality to populate the prompt folder with predefined templates.
//
// Usage:
//
//	codegpt prompt [flags]
//
// Flags:
//
//	-l, --load    load default prompt data into the specified folder (required to execute)
//
// This command will:
// 1. Check if the load flag is enabled
// 2. Get the prompt folder path from configuration
// 3. Ask for user confirmation before proceeding with data loading
// 4. Save all default prompt templates to the specified folder
//
// The command requires explicit confirmation from the user as it may overwrite existing prompt data.
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Load default prompt templates",
	Long: `Load default prompt templates into your configuration directory.
	
This command allows you to initialize or update your prompt templates with the default set provided by CodeGPT.
When executed with the --load flag, it will copy all standard templates to your configured prompt folder.`,
	Example: "  codegpt prompt --load",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !loadPromptData {
			return nil
		}

		folder := viper.GetString("prompt.folder")

		color.Yellow("Prompt folder: %s", folder)
		confirm, err := confirmation.New("Do you want to load the default prompt data? This will overwrite your existing data.", confirmation.No).RunPrompt()
		if err != nil || !confirm {
			return err
		}

		for _, key := range defaultPromptDataKeys {
			if err := savePromptData(folder, key); err != nil {
				return err
			}
		}
		return nil
	},
}

func savePromptData(folder, key string) error {
	// load default prompt data
	out, err := prompt.GetRawData(key)
	if err != nil {
		return err
	}

	// save out to file
	target := path.Join(folder, key)
	if err := os.WriteFile(target, out, 0o600); err != nil {
		return err
	}
	color.Cyan("save %s to %s", key, target)
	return nil
}
