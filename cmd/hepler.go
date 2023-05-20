package cmd

import (
	"errors"
	"fmt"

	"github.com/appleboy/CodeGPT/openai"
	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/util"
	"github.com/appleboy/com/file"

	"github.com/spf13/viper"
)

func check() error {
	// Check if the Git command is available on the system's PATH
	if !util.IsCommandAvailable("git") {
		return errors.New("Git command not found on your system's PATH. Please install Git and try again.")
	}

	// Update Viper configuration values based on the CLI flags
	if diffUnified != 3 {
		viper.Set("git.diff_unified", diffUnified)
	}

	if len(excludeList) > 0 {
		viper.Set("git.exclude_list", excludeList)
	}

	if prompt.GetLanguage(commitLang) != prompt.DefaultLanguage {
		viper.Set("output.lang", commitLang)
	}

	if openai.GetModel(commitModel) != openai.DefaultModel {
		viper.Set("openai.model", commitModel)
	}

	if httpsProxy != "" {
		viper.Set("openai.proxy", httpsProxy)
	}

	if socksProxy != "" {
		viper.Set("openai.socks", socksProxy)
	}

	if maxTokens != 0 {
		viper.Set("openai.max_tokens", maxTokens)
	}

	if templateFile != "" {
		viper.Set("git.template_file", templateFile)
	}

	if templateString != "" {
		viper.Set("git.template_string", templateString)
	}

	// Check if the template file specified in the configuration exists
	templateFile := viper.GetString("git.template_file")
	if templateFile != "" && !file.IsFile(templateFile) {
		return fmt.Errorf("template file not found: %s", templateFile)
	}

	if templateVarsFile != "" && !file.IsFile(templateVarsFile) {
		return fmt.Errorf("template variables file not found: %s", templateVarsFile)
	}

	return nil
}
