package cmd

import (
	"fmt"
	"html"
	"os"
	"path"
	"strings"
	"time"

	"github.com/appleboy/CodeGPT/core"
	"github.com/appleboy/CodeGPT/git"
	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/confirmation"
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
	noConfirm      = false
)

func init() {
	commitCmd.PersistentFlags().StringP("file", "f", "", "specify output file for commit message")
	commitCmd.PersistentFlags().BoolVar(&preview, "preview", false, "preview commit message before committing")
	commitCmd.PersistentFlags().IntVar(&diffUnified, "diff_unified", 3,
		"generate diffs with <n> lines of context (default: 3)")
	commitCmd.PersistentFlags().StringVar(&commitModel, "model", "gpt-4o", "specify which OpenAI model to use for generation")
	commitCmd.PersistentFlags().StringVar(&commitLang, "lang", "en", "set output language for the commit message (default: English)")
	commitCmd.PersistentFlags().StringSliceVar(&excludeList, "exclude_list", []string{},
		"specify files to exclude from git diff")
	commitCmd.PersistentFlags().StringVar(&httpsProxy, "proxy", "", "set HTTP proxy URL")
	commitCmd.PersistentFlags().StringVar(&socksProxy, "socks", "", "set SOCKS proxy URL")
	commitCmd.PersistentFlags().StringVar(&templateFile, "template_file", "", "provide template file for commit message format")
	commitCmd.PersistentFlags().StringVar(&templateString, "template_string", "", "provide inline template string for commit message format")
	commitCmd.PersistentFlags().StringSliceVar(&templateVars, "template_vars", []string{}, "define custom variables for templates")
	commitCmd.PersistentFlags().StringVar(&templateVarsFile, "template_vars_file", "", "specify file containing template variables")
	commitCmd.PersistentFlags().BoolVar(&commitAmend, "amend", false,
		"amend the previous commit instead of creating a new one")
	commitCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", defaultTimeout, "set API request timeout duration")
	commitCmd.PersistentFlags().BoolVar(&promptOnly, "prompt_only", false,
		"display the prompt without sending to OpenAI")
	commitCmd.PersistentFlags().BoolVar(&noConfirm, "no_confirm", false,
		"skip all confirmation prompts")
	_ = viper.BindPFlag("output.file", commitCmd.PersistentFlags().Lookup("file"))
}

// commitCmd represents the commit command.
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Automatically generate commit message",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := check(cmd.Context()); err != nil {
			return err
		}

		g := git.New(
			git.WithDiffUnified(viper.GetInt("git.diff_unified")),
			git.WithExcludeList(viper.GetStringSlice("git.exclude_list")),
			git.WithEnableAmend(commitAmend),
		)
		diff, err := g.DiffFiles(cmd.Context())
		if err != nil {
			return err
		}

		// Update the OpenAI client request timeout if the timeout value is greater than the default openai.timeout
		if timeout > viper.GetDuration("openai.timeout") ||
			timeout != defaultTimeout {
			viper.Set("openai.timeout", timeout)
		}

		// Check provider
		provider := core.Platform(viper.GetString("openai.provider"))
		client, err := GetClient(cmd.Context(), provider)
		if err != nil && !promptOnly {
			return err
		}

		currentModel := viper.GetString("openai.model")
		color.Green("Summarizing commit message using " + currentModel + " model")

		data := util.Data{}
		// Add template variables
		if vars := util.ConvertToMap(templateVars); len(vars) > 0 {
			for k, v := range vars {
				data[k] = v
			}
		}

		// Add template variables from file
		if templateVarsFile != "" {
			allENV, err := godotenv.Read(templateVarsFile)
			if err != nil {
				return err
			}
			for k, v := range allENV {
				data[k] = v
			}
		}

		// Get code review message from diff data
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

			// Determine if the user wants to use the prompt only
			if promptOnly {
				color.Yellow("====================Prompt========================")
				color.Yellow("\n" + strings.TrimSpace(out) + "\n\n")
				color.Yellow("==================================================")
				return nil
			}

			// Get summarized comment from diff data
			color.Cyan("Summarizing git diff...")
			resp, err := client.Completion(cmd.Context(), out)
			if err != nil {
				return err
			}
			data[prompt.SummarizeMessageKey] = strings.TrimSpace(resp.Content)
			color.Magenta(resp.Usage.String())
		}

		// Get summarized title from diff data
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

			// Generate title for pull request with retry if empty
			color.Cyan("Generating title for pull request...")
			const maxRetries = 3
			const retryDelay = 500 * time.Millisecond

			var summarizeTitle string
			var resp *core.Response

			for attempt := 1; attempt <= maxRetries; attempt++ {
				resp, err = client.Completion(cmd.Context(), out)
				if err != nil {
					return err
				}

				summarizeTitle = strings.TrimSpace(resp.Content)
				color.Magenta(resp.Usage.String())

				if len(summarizeTitle) > 0 {
					break
				}

				if attempt < maxRetries {
					color.Cyan("Empty title response, retrying (%d/%d)...", attempt, maxRetries)
					time.Sleep(retryDelay)
				}
			}

			if len(summarizeTitle) == 0 {
				return fmt.Errorf("failed to get valid title after %d attempts", maxRetries)
			}

			// Lowercase the first character of first word of the commit message and remove the trailing period
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
			message := "Generating conventional commit prefix"
			summaryPrix := ""
			color.Cyan(message + " (Tools)")
			resp, err := client.GetSummaryPrefix(cmd.Context(), out)
			if err != nil {
				return err
			}
			summaryPrix = resp.Content

			color.Magenta(resp.Usage.String())

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

			// Translate git commit message
			color.Cyan("Translating git commit message to " + prompt.GetLanguage(viper.GetString("output.lang")))
			resp, err := client.Completion(cmd.Context(), out)
			if err != nil {
				return err
			}
			color.Magenta(resp.Usage.String())
			commitMessage = resp.Content
		}

		// Unescape HTML entities in commit message
		commitMessage = html.UnescapeString(commitMessage)
		commitMessage = strings.TrimSpace(commitMessage)

		// Output commit summary data from AI
		color.Yellow("================Commit Summary====================")
		color.Yellow("\n" + commitMessage + "\n\n")
		color.Yellow("==================================================")

		outputFile := viper.GetString("output.file")
		if outputFile == "" {
			out, err := g.GitDir(cmd.Context())
			if err != nil {
				return err
			}
			outputFile = path.Join(strings.TrimSpace(out), "COMMIT_EDITMSG")
		}
		color.Cyan("Writing commit message to " + outputFile)
		// Write commit message to git staging file
		err = os.WriteFile(outputFile, []byte(commitMessage), 0o600)
		if err != nil {
			return err
		}

		// Handle preview: if preview and noConfirm, or preview prompt declined, then exit early
		if preview {
			if noConfirm {
				return nil
			}
			if ready, err := confirmation.New("Commit this preview summary?", confirmation.Yes).RunPrompt(); err != nil || !ready {
				if err != nil {
					return err
				}
				return nil
			}
		}

		// Handle commit message change prompt when confirmation is enabled
		if !noConfirm {
			if change, err := confirmation.New("Do you want to modify the commit message?", confirmation.No).RunPrompt(); err != nil {
				return err
			} else if change {
				m := initialPrompt(commitMessage)
				p := tea.NewProgram(m, tea.WithContext(cmd.Context()))
				if _, err := p.Run(); err != nil {
					return err
				}
				p.Wait()
				commitMessage = m.textarea.Value()
			}
		}

		// Commit changes to repository
		color.Cyan("Recording changes to the repository")
		output, err := g.Commit(cmd.Context(), commitMessage)
		if err != nil {
			return err
		}
		color.Yellow(output)
		return nil
	},
}
