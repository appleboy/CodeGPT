package prompt

import (
	"embed"
	"log"

	"github.com/appleboy/CodeGPT/util"
)

//go:embed templates/*
var templatesFS embed.FS

// Template file names
const (
	CodeReviewTemplate         = "code_review_file_diff.tmpl"
	SummarizeFileDiffTemplate  = "summarize_file_diff.tmpl"
	SummarizeTitleTemplate     = "summarize_title.tmpl"
	ConventionalCommitTemplate = "conventional_commit.tmpl"
	TranslationTemplate        = "translation.tmpl"
	SummarizePrefixKey         = "summarize_prefix"
	SummarizeTitleKey          = "summarize_title"
	SummarizeMessageKey        = "summarize_message"
)

// Initializes the prompt package by loading the templates from the embedded file system.
func init() {
	if err := util.LoadTemplates(templatesFS); err != nil {
		log.Fatal(err)
	}
}
