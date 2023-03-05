package prompt

import (
	"embed"
	"log"

	"github.com/appleboy/CodeGPT/util"
)

//go:embed templates/*
var files embed.FS

const (
	SummarizeFileDiffTemplate = "summarize_file_diff.tmpl"
	SummarizeTitleTemplate    = "summarize_title.tmpl"
	TranslationTemplate       = "translation.tmpl"
)

func init() {
	if err := util.LoadTemplates(files); err != nil {
		log.Fatal(err)
	}
}
