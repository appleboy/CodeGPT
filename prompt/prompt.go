package prompt

import (
	"embed"
	"log"

	"github.com/appleboy/CodeGPT/util"
)

//go:embed templates/*
var files embed.FS

const (
	CodeReviewTemplate         = "code_review_file_diff.tmpl"
	SummarizeFileDiffTemplate  = "summarize_file_diff.tmpl"
	SummarizeTitleTemplate     = "summarize_title.tmpl"
	ConventionalCommitTemplate = "conventional_commit.tmpl"
	TranslationTemplate        = "translation.tmpl"
)

func init() {
	if err := util.LoadTemplates(files); err != nil {
		log.Fatal(err)
	}
}
