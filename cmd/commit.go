package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/appleboy/CodeGPT/git"
	"github.com/appleboy/CodeGPT/openai"
	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/util"

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
	Run: func(cmd *cobra.Command, args []string) {
		var message string

		// check git command exist
		if !util.IsCommandAvailable("git") {
			log.Fatal("To use CodeGPT, you must have git on your PATH")
		}

		diff, err := git.Diff()
		if err != nil {
			log.Fatal(err)
		}

		client, err := openai.New(
			viper.GetString("openai.api_key"),
			viper.GetString("openai.model"),
			viper.GetString("openai.org_id"),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Get summarize comment from diff datas
		out, err := util.GetTemplate(
			prompt.SummarizeFileDiffTemplate,
			util.Data{
				"file_diffs": diff,
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		summarizeDiff, err := client.Completion(cmd.Context(), out)
		if err != nil {
			log.Fatal(err)
		}

		out, err = util.GetTemplate(
			prompt.SummarizeTitleTemplate,
			util.Data{
				"summary_points": summarizeDiff,
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		summarizeTitle, err := client.Completion(cmd.Context(), out)
		if err != nil {
			log.Fatal(err)
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
				log.Fatal(err)
			}

			summarize, err := client.Completion(cmd.Context(), out)
			if err != nil {
				log.Fatal(err)
			}
			message = summarize
		} else {
			message = strings.TrimSpace(summarizeTitle) + "\n\n" + strings.TrimSpace(summarizeDiff)
		}

		err = os.WriteFile(viper.GetString("output.file"), []byte(message), 0o644)
		if err != nil {
			log.Fatal(err)
		}
	},
}
