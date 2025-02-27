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
		"load default prompt data")
}

var defaultPromptDataKeys = []string{
	prompt.CodeReviewTemplate,
	prompt.SummarizeFileDiffTemplate,
	prompt.SummarizeTitleTemplate,
	prompt.ConventionalCommitTemplate,
}

// promptCmd represents the command to load default prompt data.
// It uses the "prompt" keyword and provides a short description: "load default prompt data".
// The command executes the RunE function which checks if the loadPromptData flag is set.
// If set, it prompts the user for confirmation to load the default prompt data, which will overwrite existing data.
// Upon confirmation, it retrieves the prompt folder path from the configuration and saves the default prompt data keys to the specified folder.
// If any error occurs during the process, it returns the error.
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "load default prompt data",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !loadPromptData {
			return nil
		}

		confirm, err := confirmation.New("Do you want to load default prompt data, will overwrite your data", confirmation.No).RunPrompt()
		if err != nil || !confirm {
			return err
		}

		folder := viper.GetString("prompt_folder")
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
