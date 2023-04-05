package cmd

import (
	"strconv"
	"strings"

	"github.com/appleboy/CodeGPT/git"
	"github.com/appleboy/CodeGPT/openai"
	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// The maximum number of tokens to generate in the chat completion.
// The total length of input tokens and generated tokens is limited by the model's context length.
var maxTokens int

func init() {
	reviewCmd.Flags().IntVar(&diffUnified, "diff_unified", 3, "generate diffs with <n> lines of context, default is 3")
	reviewCmd.Flags().IntVar(&maxTokens, "max_tokens", 300, "the maximum number of tokens to generate in the chat completion.")
	reviewCmd.Flags().StringVar(&commitModel, "model", "gpt-3.5-turbo", "select openai model")
	reviewCmd.Flags().StringVar(&commitLang, "lang", "en", "summarizing language uses English by default")
	reviewCmd.Flags().StringSliceVar(&excludeList, "exclude_list", []string{}, "exclude file from git diff command")
	reviewCmd.Flags().BoolVar(&commitAmend, "amend", false, "replace the tip of the current branch by creating a new commit.")
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Auto review code changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := check(); err != nil {
			return err
		}

		g := git.New(
			git.WithDiffUnified(viper.GetInt("git.diff_unified")),
			git.WithExcludeList(viper.GetStringSlice("git.exclude_list")),
			git.WithEnableAmend(commitAmend),
		)
		diff, err := g.DiffFiles()
		if err != nil {
			return err
		}

		color.Green("Code review your changes using " + viper.GetString("openai.model") + " model")
		client, err := openai.New(
			openai.WithToken(viper.GetString("openai.api_key")),
			openai.WithModel(viper.GetString("openai.model")),
			openai.WithOrgID(viper.GetString("openai.org_id")),
			openai.WithProxyURL(viper.GetString("openai.proxy")),
			openai.WithSocksURL(viper.GetString("openai.socks")),
			openai.WithBaseURL(viper.GetString("openai.base_url")),
			openai.WithTimeout(viper.GetDuration("openai.timeout")),
			openai.WithMaxTokens(viper.GetInt("openai.max_tokens")),
			openai.WithTemperature(float32(viper.GetFloat64("openai.temperature"))),
		)
		if err != nil {
			return err
		}

		out, err := util.GetTemplateByString(
			prompt.CodeReviewTemplate,
			util.Data{
				"file_diffs": diff,
			},
		)
		if err != nil {
			return err
		}

		// Get summarize comment from diff datas
		color.Cyan("We are trying to review code changes")
		resp, err := client.Completion(cmd.Context(), out)
		if err != nil {
			return err
		}
		summarizeMessage := resp.Content
		color.Magenta("PromptTokens: " + strconv.Itoa(resp.Usage.PromptTokens) +
			", CompletionTokens: " + strconv.Itoa(resp.Usage.CompletionTokens) +
			", TotalTokens: " + strconv.Itoa(resp.Usage.TotalTokens),
		)

		if prompt.GetLanguage(viper.GetString("output.lang")) != prompt.DefaultLanguage {
			out, err = util.GetTemplateByString(
				prompt.TranslationTemplate,
				util.Data{
					"output_language": prompt.GetLanguage(viper.GetString("output.lang")),
					"output_message":  summarizeMessage,
				},
			)
			if err != nil {
				return err
			}

			// translate a git commit message
			color.Cyan("We are trying to translate code review to " + prompt.GetLanguage(viper.GetString("output.lang")) + " language")
			resp, err := client.Completion(cmd.Context(), out)
			if err != nil {
				return err
			}
			color.Magenta("PromptTokens: " + strconv.Itoa(resp.Usage.PromptTokens) +
				", CompletionTokens: " + strconv.Itoa(resp.Usage.CompletionTokens) +
				", TotalTokens: " + strconv.Itoa(resp.Usage.TotalTokens),
			)
			summarizeMessage = resp.Content
		}

		// Output core review summary
		color.Yellow("================Review Summary====================")
		color.Yellow("\n" + strings.TrimSpace(summarizeMessage) + "\n\n")
		color.Yellow("==================================================")

		return nil
	},
}
