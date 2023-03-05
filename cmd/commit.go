package cmd

import (
	"fmt"
	"log"

	"github.com/appleboy/CodeGPT/git"
	"github.com/appleboy/CodeGPT/openai"
	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	commitCmd.PersistentFlags().StringP("model", "m", "gpt-3.5-turbo", "openai model")
	commitCmd.PersistentFlags().StringP("lang", "l", "en", "summarizing language uses English by default")
	_ = viper.BindPFlag("openai.model", commitCmd.PersistentFlags().Lookup("model"))
	_ = viper.BindPFlag("output.lang", commitCmd.PersistentFlags().Lookup("lang"))
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Auto generate commit message",
	Run: func(cmd *cobra.Command, args []string) {
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
		out, err := prompt.GetTemplate(
			prompt.SummarizeFileDiffTemplate,
			prompt.Data{
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

		out, err = prompt.GetTemplate(
			prompt.SummarizeTitleTemplate,
			prompt.Data{
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

		out, err = prompt.GetTemplate(
			prompt.TranslationTemplate,
			prompt.Data{
				"output_language": prompt.GetLanguage(viper.GetString("output.lang")),
				"commit_title":    summarizeTitle,
				"commit_message":  summarizeDiff,
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		message, err := client.Completion(cmd.Context(), out)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(message)
	},
}
