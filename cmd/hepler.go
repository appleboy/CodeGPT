package cmd

import (
	"errors"

	"github.com/appleboy/CodeGPT/openai"
	"github.com/appleboy/CodeGPT/prompt"
	"github.com/appleboy/CodeGPT/util"
	"github.com/appleboy/com/file"

	"github.com/spf13/viper"
)

func check() error {
	// check git command exist
	if !util.IsCommandAvailable("git") {
		return errors.New("Git command not found on your system's PATH. Please install Git and try again.")
	}

	if diffUnified != 3 {
		viper.Set("git.diff_unified", diffUnified)
	}

	if len(excludeList) > 0 {
		viper.Set("git.exclude_list", excludeList)
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

	if maxTokens != 0 {
		viper.Set("openai.max_tokens", maxTokens)
	}

	if templateFile != "" {
		viper.Set("git.template_file", templateFile)
	}

	if templateString != "" {
		viper.Set("git.template_string", templateString)
	}

	if viper.GetString("git.template_file") != "" {
		if !file.IsFile(viper.GetString("git.template_file")) {
			return errors.New("template file not found: " + viper.GetString("git.template_file"))
		}
	}

	return nil
}
