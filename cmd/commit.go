package cmd

import (
	"errors"
	"html"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/git"
	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/util"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	commitLang  string
	commitModel string

	preview        bool
	diffUnified    int
	excludeList    []string
	httpsProxy     string
	socksProxy     string
	templateFile   string
	templateString string
	commitAmend    bool
	timeout        time.Duration
	promptOnly     bool

	templateVars     []string
	templateVarsFile string

	defaultTimeout = 30 * time.Second
)

func init() {
	commitCmd.PersistentFlags().StringP("file", "f", "", "commit message file")
	commitCmd.PersistentFlags().BoolVar(&preview, "preview", false, "preview commit message")
	commitCmd.PersistentFlags().IntVar(&diffUnified, "diff_unified", 3,
		"generate diffs with <n> lines of context, default is 3")
	commitCmd.PersistentFlags().StringVar(&commitModel, "model", "gpt-3.5-turbo", "select openai model")
	commitCmd.PersistentFlags().StringVar(&commitLang, "lang", "en", "summarizing language uses English by default")
	commitCmd.PersistentFlags().StringSliceVar(&excludeList, "exclude_list", []string{},
		"exclude file from git diff command")
	commitCmd.PersistentFlags().StringVar(&httpsProxy, "proxy", "", "http proxy")
	commitCmd.PersistentFlags().StringVar(&socksProxy, "socks", "", "socks proxy")
	commitCmd.PersistentFlags().StringVar(&templateFile, "template_file", "", "git commit message file")
	commitCmd.PersistentFlags().StringVar(&templateString, "template_string", "", "git commit message string")
	commitCmd.PersistentFlags().StringSliceVar(&templateVars, "template_vars", []string{}, "template variables")
	commitCmd.PersistentFlags().StringVar(&templateVarsFile, "template_vars_file", "", "template variables file")
	commitCmd.PersistentFlags().BoolVar(&commitAmend, "amend", false,
		"replace the tip of the current branch by creating a new commit.")
	commitCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", defaultTimeout, "request timeout")
	commitCmd.PersistentFlags().BoolVar(&promptOnly, "prompt_only", false,
		"show prompt only, don't send request to openai")
	_ = viper.BindPFlag("output.file", commitCmd.PersistentFlags().Lookup("file"))
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Auto generate commit message",
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
		if !provider.IsValid() {
			return errors.New("invalid provider")
		}

		var client core.Generative
		switch provider {
		case core.Gemini:
			// TODO: implement Gemini
		case core.OpenAI, core.Azure:
			client, err = NewOpenAI()
		}
		if err != nil && !promptOnly {
			return err
		}

		currentModel := viper.GetString("openai.model")
		color.Green("Summarize the commit message use " + currentModel + " model")

		data := util.Data{}
		// add template vars
		if vars := util.ConvertToMap(templateVars); len(vars) > 0 {
			for k, v := range vars {
				data[k] = v
			}
		}

		// add template vars from file
		if templateVarsFile != "" {
			allENV, err := godotenv.Read(templateVarsFile)
			if err != nil {
				return err
			}
			for k, v := range allENV {
				data[k] = v
			}
		}

		// Get code review message from diff datas
		if _, ok := data[prompt.SummarizeMessageKey]; !ok {
			out, err := util.GetTemplateByString(
				prompt.SummarizeFileDiffTemplate,
				util.Data{
					"file_diffs": diff,
				},
			)
			if err != nil {
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
			color.Cyan("We are trying to summarize a git diff")
			resp, err := client.Completion(cmd.Context(), out)
			if err != nil {
				return err
			}
			data[prompt.SummarizeMessageKey] = strings.TrimSpace(resp.Content)
			color.Magenta("PromptTokens: " + strconv.Itoa(resp.Usage.PromptTokens) +
				", CompletionTokens: " + strconv.Itoa(resp.Usage.CompletionTokens) +
				", TotalTokens: " + strconv.Itoa(resp.Usage.TotalTokens),
			)
		}

		// Get summarize title from diff datas
		if _, ok := data[prompt.SummarizeTitleKey]; !ok {
			out, err := util.GetTemplateByString(
				prompt.SummarizeTitleTemplate,
				util.Data{
					"summary_points": data[prompt.SummarizeMessageKey],
				},
			)
			if err != nil {
				return err
			}

			// Get summarize title from diff datas
			color.Cyan("We are trying to summarize a title for pull request")
			resp, err := client.Completion(cmd.Context(), out)
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
			data[prompt.SummarizeTitleKey] = strings.TrimSpace(summarizeTitle)
		}

		if _, ok := data[prompt.SummarizePrefixKey]; !ok {
			out, err := util.GetTemplateByString(
				prompt.ConventionalCommitTemplate,
				util.Data{
					"summary_points": data[prompt.SummarizeMessageKey],
				},
			)
			if err != nil {
				return err
			}
			message := "We are trying to get conventional commit prefix"
			summaryPrix := ""
			color.Cyan(message + " (Tools)")
			resp, err := client.GetSummaryPrefix(cmd.Context(), out)
			if err != nil {
				return err
			}
			summaryPrix = resp.Content

			color.Magenta("PromptTokens: " + strconv.Itoa(resp.Usage.PromptTokens) +
				", CompletionTokens: " + strconv.Itoa(resp.Usage.CompletionTokens) +
				", TotalTokens: " + strconv.Itoa(resp.Usage.TotalTokens),
			)

			data[prompt.SummarizePrefixKey] = summaryPrix
		}

		var commitMessage string
		if viper.GetString("git.template_file") != "" {
			format, err := os.ReadFile(viper.GetString("git.template_file"))
			if err != nil {
				return err
			}
			commitMessage, err = util.NewTemplateByString(
				string(format),
				data,
			)
			if err != nil {
				return err
			}
		} else if viper.GetString("git.template_string") != "" {
			commitMessage, err = util.NewTemplateByString(
				viper.GetString("git.template_string"),
				data,
			)
			if err != nil {
				return err
			}
		} else {
			commitMessage, err = util.GetTemplateByString(
				git.CommitMessageTemplate,
				data,
			)
			if err != nil {
				return err
			}
		}

		if prompt.GetLanguage(viper.GetString("output.lang")) != prompt.DefaultLanguage {
			out, err := util.GetTemplateByString(
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

		// unescape html entities in commit message
		commitMessage = html.UnescapeString(commitMessage)

		// Output commit summary data from AI
		color.Yellow("================Commit Summary====================")
		color.Yellow("\n" + strings.TrimSpace(commitMessage) + "\n\n")
		color.Yellow("==================================================")

		outputFile := viper.GetString("output.file")
		if outputFile == "" {
			out, err := g.GitDir()
			if err != nil {
				return err
			}
			outputFile = path.Join(strings.TrimSpace(out), "COMMIT_EDITMSG")
		}
		color.Cyan("Write the commit message to " + outputFile + " file")
		// write commit message to git staging file
		err = os.WriteFile(outputFile, []byte(commitMessage), 0o600)
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
