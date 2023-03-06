package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/appleboy/CodeGPT/git"
	"github.com/appleboy/CodeGPT/openai"
	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	commitCmd.PersistentFlags().StringP("file", "f", ".git/COMMIT_EDITMSG", "commit message file")
	commitCmd.PersistentFlags().StringP("model", "m", "gpt-3.5-turbo", "openai model")
	commitCmd.PersistentFlags().StringP("lang", "l", "en", "summarizing language uses English by default")
	_ = viper.BindPFlag("openai.model", commitCmd.PersistentFlags().Lookup("model"))
	_ = viper.BindPFlag("output.lang", commitCmd.PersistentFlags().Lookup("lang"))
	_ = viper.BindPFlag("output.file", commitCmd.PersistentFlags().Lookup("file"))
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Auto generate commit message",
	RunE: func(cmd *cobra.Command, args []string) error {
		var message string

		// check git command exist
		if !util.IsCommandAvailable("git") {
			return errors.New("To use CodeGPT, you must have git on your PATH")
		}

		diff, err := git.Diff()
		if err != nil {
			return err
		}

		color.Green("Summarize the commit message use " + viper.GetString("openai.model") + " model")

		client, err := openai.New(
			viper.GetString("openai.api_key"),
			viper.GetString("openai.model"),
			viper.GetString("openai.org_id"),
			viper.GetString("openai.proxy"),
		)
		if err != nil {
			return err
		}

		// Get summarize comment from diff datas
		out, err := util.GetTemplate(
			prompt.SummarizeFileDiffTemplate,
			util.Data{
				"file_diffs": diff,
			},
		)
		if err != nil {
			return err
		}

		color.Cyan("We are trying to summarize a git diff")
		summarizeDiff, err := client.Completion(cmd.Context(), out)
		if err != nil {
			return err
		}

		out, err = util.GetTemplate(
			prompt.SummarizeTitleTemplate,
			util.Data{
				"summary_points": summarizeDiff,
			},
		)
		if err != nil {
			return err
		}

		color.Cyan("We are trying to summarize a title for pull request")
		summarizeTitle, err := client.Completion(cmd.Context(), out)
		if err != nil {
			return err
		}

		if prompt.GetLanguage(viper.GetString("output.lang")) != prompt.DefaultLanguage {
			out, err = util.GetTemplate(
				prompt.TranslationTemplate,
				util.Data{
					"output_language": prompt.GetLanguage(viper.GetString("output.lang")),
					"commit_title":    summarizeTitle,
					"commit_message":  summarizeDiff,
				},
			)
			if err != nil {
				return err
			}

			color.Cyan("We are trying to translate a git commit message to " + prompt.GetLanguage(viper.GetString("output.lang")) + "language")
			summarize, err := client.Completion(cmd.Context(), out)
			if err != nil {
				return err
			}
			message = summarize
		} else {
			message = strings.TrimSpace(summarizeTitle) + "\n\n" + strings.TrimSpace(summarizeDiff)
		}
		color.Yellow("================Commit Summary====================")
		color.Yellow("\n" + message + "\n\n")
		color.Yellow("==================================================")
		color.Cyan("Write the commit message to " + viper.GetString("output.file") + " file")
		err = os.WriteFile(viper.GetString("output.file"), []byte(message), 0o644)
		if err != nil {
			return err
		}
		return nil
	},
}
