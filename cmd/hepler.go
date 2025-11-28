package cmd

import (
	"errors"
	"fmt"

	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/provider/openai"
	"github.com/appleboy/CodeGPT/util"
	"github.com/appleboy/com/file"
	"github.com/spf13/viper"
)

func check() error {
	// Check if the Git command is available on the system's PATH
	if !util.IsCommandAvailable("git") {
		return errors.New("git command not found in your system's PATH. Please install Git and try again")
	}

	// Check if current working directory is inside a Git repository.
	if !util.IsGitRepo() {
		return errors.New("not a git repository")
	}

	// Apply configuration values from CLI flags to Viper
	if diffUnified != 3 {
		viper.Set("git.diff_unified", diffUnified)
	}

	if len(excludeList) > 0 {
		viper.Set("git.exclude_list", excludeList)
	}

	if prompt.GetLanguage(commitLang) != prompt.DefaultLanguage {
		viper.Set("output.lang", commitLang)
	}

	if commitModel != openai.DefaultModel {
		viper.Set("openai.model", commitModel)
	}

	if httpsProxy != "" {
		viper.Set("openai.proxy", httpsProxy)
	}

	if socksProxy != "" {
		viper.Set("openai.socks", socksProxy)
	}

	if maxTokens != 300 {
		viper.Set("openai.max_tokens", maxTokens)
	}

	if templateFile != "" {
		viper.Set("git.template_file", templateFile)
	}

	if templateString != "" {
		viper.Set("git.template_string", templateString)
	}

	// Verify template file existence
	cfgTemplateFile := viper.GetString("git.template_file")
	if cfgTemplateFile != "" {
		exists, err := file.IsFile(cfgTemplateFile)
		if err != nil {
			return fmt.Errorf("failed to check template file: %v", err)
		}
		if !exists {
			return fmt.Errorf("template file not found at: %s", cfgTemplateFile)
		}
	}

	if templateVarsFile != "" {
		exists, err := file.IsFile(templateVarsFile)
		if err != nil {
			return fmt.Errorf("failed to check template variables file: %v", err)
		}
		if !exists {
			return fmt.Errorf("template variables file not found at: %s", templateVarsFile)
		}
	}

	// Load custom prompts from configured directory
	promptFolder := viper.GetString("prompt.folder")
	if promptFolder != "" {
		if err := util.LoadTemplatesFromDir(promptFolder); err != nil {
			return fmt.Errorf("failed to load custom prompt templates: %s", err)
		}
	}

	return nil
}
