package cmd

import (
	"strconv"
	"strings"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/git"
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
	reviewCmd.PersistentFlags().IntVar(&diffUnified, "diff_unified", 3,
		"generate diffs with <n> lines of context, default is 3")
	reviewCmd.PersistentFlags().IntVar(&maxTokens, "max_tokens", 300,
		"the maximum number of tokens to generate in the chat completion.")
	reviewCmd.PersistentFlags().StringVar(&commitModel, "model", "gpt-3.5-turbo", "select openai model")
	reviewCmd.PersistentFlags().StringVar(&commitLang, "lang", "en", "summarizing language uses English by default")
	reviewCmd.PersistentFlags().StringSliceVar(&excludeList, "exclude_list", []string{}, "exclude file from git diff command")
	reviewCmd.PersistentFlags().BoolVar(&commitAmend, "amend", false,
		"replace the tip of the current branch by creating a new commit.")
	reviewCmd.PersistentFlags().BoolVar(&promptOnly, "prompt_only", false,
		"show prompt only, don't send request to openai")
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

		// Update the OpenAI client request timeout if the timeout value is greater than the default openai.timeout
		if timeout > viper.GetDuration("openai.timeout") ||
			timeout != defaultTimeout {
			viper.Set("openai.timeout", timeout)
		}

		// check provider
		provider := core.Platform(viper.GetString("openai.provider"))
		client, err := GetClient(provider)
		if err != nil {
			return err
		}

		currentModel := viper.GetString("openai.model")
		color.Green("Code review your changes using " + currentModel + " model")

		out, err := util.GetTemplateByString(
			prompt.CodeReviewTemplate,
			util.Data{
				"file_diffs": diff,
			},
		)
		if err != nil && !promptOnly {
			return err
		}

		// determine if the user wants to use the prompt only
		if promptOnly {
			color.Yellow("====================Prompt========================")
			color.Yellow("\n" + strings.TrimSpace(out) + "\n\n")
			color.Yellow("==================================================")
			return nil
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
			color.Cyan("we are trying to translate code review to " +
				prompt.GetLanguage(viper.GetString("output.lang")) + " language")
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
