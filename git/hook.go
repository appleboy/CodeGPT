// Package git provides functionality for working with git repositories.
//
// This package includes constants defining template names for git hooks and
// initializes the templates from embedded files on package initialization.
//
// The package embeds files from a 'templates/' directory which contain hook templates
// that can be used for git operations such as preparing commit messages.
package git

import (
	"embed"
	"log"

	"github.com/appleboy/CodeGPT/util"
)

//go:embed templates/*
var files embed.FS

const (
	// HookPrepareCommitMessageTemplate is the template for the prepare-commit-msg hook
	HookPrepareCommitMessageTemplate = "prepare-commit-msg"
	// CommitMessageTemplate is the template for the commit message
	CommitMessageTemplate = "commit-msg.tmpl"
)

// init initializes the Git hook templates by loading them from embedded files.
// If there's an error loading the templates, the function logs a fatal error and terminates the program.
func init() { //nolint:gochecknoinits
	if err := util.LoadTemplates(files); err != nil {
		log.Fatal(err)
	}
}
