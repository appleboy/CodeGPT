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

var (
	commitLang  string
	commitModel string
)

func init() {
	commitCmd.PersistentFlags().StringP("file", "f", ".git/COMMIT_EDITMSG", "commit message file")
	commitCmd.PersistentFlags().StringVar(&commitModel, "model", "gpt-3.5-turbo", "select openai model")
	commitCmd.PersistentFlags().StringVar(&commitLang, "lang", "en", "summarizing language uses English by default")
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

		g := git.New()
		diff, err := g.DiffFiles()
		if err != nil {
			return err
		}

		// check default language
		if prompt.GetLanguage(commitLang) != prompt.DefaultLanguage {
			viper.Set("output.lang", commitLang)
		}

		// check default model
		if openai.GetModel(commitModel) != openai.DefaultModel {
			viper.Set("openai.model", commitModel)
		}

		color.Green("Summarize the commit message use " + viper.GetString("openai.model") + " model")
		client, err := openai.New(
			openai.WithToken(viper.GetString("openai.api_key")),
			openai.WithModel(viper.GetString("openai.model")),
			openai.WithOrgID(viper.GetString("openai.org_id")),
			openai.WithProxyURL(viper.GetString("openai.proxy")),
		)
		if err != nil {
			return err
		}

		out, err := util.GetTemplate(
			prompt.SummarizeFileDiffTemplate,
			util.Data{
				"file_diffs": diff,
			},
		)
		if err != nil {
			return err
		}

		// Get summarize comment from diff datas
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

		// Get summarize title from diff datas
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

			// translate a git commit message
			color.Cyan("We are trying to translate a git commit message to " + prompt.GetLanguage(viper.GetString("output.lang")) + " language")
			summarize, err := client.Completion(cmd.Context(), out)
			if err != nil {
				return err
			}
			message = summarize
		} else {
			message = strings.TrimSpace(summarizeTitle) + "\n\n" + strings.TrimSpace(summarizeDiff)
		}

		// Output commit summary data from AI
		color.Yellow("================Commit Summary====================")
		color.Yellow("\n" + message + "\n\n")
		color.Yellow("==================================================")
		color.Cyan("Write the commit message to " + viper.GetString("output.file") + " file")

		// write commit message to git staging file
		err = os.WriteFile(viper.GetString("output.file"), []byte(message), 0o644)
		if err != nil {
			return err
		}
		return nil
	},
}
