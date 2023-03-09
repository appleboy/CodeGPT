package git

import (
	"embed"
	"log"

	"github.com/appleboy/CodeGPT/util"
)

//go:embed templates/*
var files embed.FS

const (
	CommitMessageTemplate = "prepare-commit-msg"
)

func init() {
	if err := util.LoadTemplates(files); err != nil {
		log.Fatal(err)
	}
}
