package cmd

import (
	"errors"
	"os"
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

var (
	commitLang  string
	commitModel string

	preview     bool
	diffUnified int
	excludeList []string
	httpsProxy  string
	socksProxy  string
)

func init() {
	commitCmd.PersistentFlags().StringP("file", "f", ".git/COMMIT_EDITMSG", "commit message file")
	commitCmd.PersistentFlags().BoolVar(&preview, "preview", false, "preview commit message")
	commitCmd.PersistentFlags().IntVar(&diffUnified, "diff_unified", 3, "generate diffs with <n> lines of context, default is 3")
	commitCmd.PersistentFlags().StringVar(&commitModel, "model", "gpt-3.5-turbo", "select openai model")
	commitCmd.PersistentFlags().StringVar(&commitLang, "lang", "en", "summarizing language uses English by default")
	commitCmd.PersistentFlags().StringSliceVar(&excludeList, "exclude_list", []string{}, "exclude file from git diff command")
	commitCmd.PersistentFlags().StringVar(&httpsProxy, "proxy", "", "http proxy")
	commitCmd.PersistentFlags().StringVar(&socksProxy, "socks", "", "socks proxy")
	_ = viper.BindPFlag("output.file", commitCmd.PersistentFlags().Lookup("file"))
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Auto generate commit message",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check git command exist
		if !util.IsCommandAvailable("git") {
			return errors.New("To use CodeGPT, you must have git on your PATH")
		}

		if diffUnified != 3 {
			viper.Set("git.diff_unified", diffUnified)
		}

		if len(excludeList) > 0 {
			viper.Set("git.exclude_list", excludeList)
		}

		g := git.New(
			git.WithDiffUnified(viper.GetInt("git.diff_unified")),
			git.WithExcludeList(viper.GetStringSlice("git.exclude_list")),
		)
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

		if httpsProxy != "" {
			viper.Set("openai.proxy", httpsProxy)
		}

		if socksProxy != "" {
			viper.Set("openai.socks", socksProxy)
		}

		color.Green("Summarize the commit message use " + viper.GetString("openai.model") + " model")
		client, err := openai.New(
			openai.WithToken(viper.GetString("openai.api_key")),
			openai.WithModel(viper.GetString("openai.model")),
			openai.WithOrgID(viper.GetString("openai.org_id")),
			openai.WithProxyURL(viper.GetString("openai.proxy")),
			openai.WithSocksURL(viper.GetString("openai.socks")),
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
		resp, err := client.Completion(cmd.Context(), out)
		if err != nil {
			return err
		}
		summarizeMessage := resp.Content
		color.Magenta("PromptTokens: " + strconv.Itoa(resp.Usage.PromptTokens) +
			", CompletionTokens: " + strconv.Itoa(resp.Usage.CompletionTokens) +
			", TotalTokens: " + strconv.Itoa(resp.Usage.TotalTokens),
		)

		out, err = util.GetTemplate(
			prompt.SummarizeTitleTemplate,
			util.Data{
				"summary_points": summarizeMessage,
			},
		)
		if err != nil {
			return err
		}

		// Get summarize title from diff datas
		color.Cyan("We are trying to summarize a title for pull request")
		resp, err = client.Completion(cmd.Context(), out)
		if err != nil {
			return err
		}
		summarizeTitle := resp.Content
		color.Magenta("PromptTokens: " + strconv.Itoa(resp.Usage.PromptTokens) +
			", CompletionTokens: " + strconv.Itoa(resp.Usage.CompletionTokens) +
			", TotalTokens: " + strconv.Itoa(resp.Usage.TotalTokens),
		)

		// lowercase the first character of first word of the commit message and remove last period
		summarizeTitle = strings.TrimRight(strings.ToLower(string(summarizeTitle[0]))+summarizeTitle[1:], ".")

		// support conventional commits
		out, err = util.GetTemplate(
			prompt.ConventionalCommitTemplate,
			util.Data{
				"summary_points": summarizeMessage,
			},
		)
		if err != nil {
			return err
		}
		color.Cyan("We are trying to get conventional commit prefix")
		resp, err = client.Completion(cmd.Context(), out)
		if err != nil {
			return err
		}
		summarizePrefix := resp.Content
		color.Magenta("PromptTokens: " + strconv.Itoa(resp.Usage.PromptTokens) +
			", CompletionTokens: " + strconv.Itoa(resp.Usage.CompletionTokens) +
			", TotalTokens: " + strconv.Itoa(resp.Usage.TotalTokens),
		)

		commitMessage, err := util.GetTemplate(
			git.CommitMessageTemplate,
			util.Data{
				"summarize_prefix":  strings.TrimSpace(summarizePrefix),
				"summarize_title":   strings.TrimSpace(summarizeTitle),
				"summarize_message": strings.TrimSpace(summarizeMessage),
			},
		)
		if err != nil {
			return err
		}

		if prompt.GetLanguage(viper.GetString("output.lang")) != prompt.DefaultLanguage {
			out, err = util.GetTemplate(
				prompt.TranslationTemplate,
				util.Data{
					"output_language": prompt.GetLanguage(viper.GetString("output.lang")),
					"output_message":  commitMessage,
				},
			)
			if err != nil {
				return err
			}

			// translate a git commit message
			color.Cyan("We are trying to translate a git commit message to " + prompt.GetLanguage(viper.GetString("output.lang")) + " language")
			resp, err := client.Completion(cmd.Context(), out)
			if err != nil {
				return err
			}
			color.Magenta("PromptTokens: " + strconv.Itoa(resp.Usage.PromptTokens) +
				", CompletionTokens: " + strconv.Itoa(resp.Usage.CompletionTokens) +
				", TotalTokens: " + strconv.Itoa(resp.Usage.TotalTokens),
			)
			commitMessage = resp.Content
		}

		// Output commit summary data from AI
		color.Yellow("================Commit Summary====================")
		color.Yellow("\n" + strings.TrimSpace(commitMessage) + "\n\n")
		color.Yellow("==================================================")

		color.Cyan("Write the commit message to " + viper.GetString("output.file") + " file")
		// write commit message to git staging file
		err = os.WriteFile(viper.GetString("output.file"), []byte(commitMessage), 0o644)
		if err != nil {
			return err
		}

		if preview {
			return nil
		}

		// git commit automatically
		color.Cyan("Git record changes to the repository")
		output, err := g.Commit(commitMessage)
		if err != nil {
			return err
		}
		color.Yellow(output)
		return nil
	},
}
