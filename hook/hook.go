package hook

import (
	"embed"
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/appleboy/CodeGPT/git"
	"github.com/appleboy/CodeGPT/util"

	"github.com/appleboy/com/file"
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

func Install() error {
	g := git.New()
	hookPath, err := g.HookPath()
	if err != nil {
		return err
	}

	target := path.Join(strings.TrimSpace(hookPath), CommitMessageTemplate)
	if file.IsFile(target) {
		return errors.New("hook file prepare-commit-msg exist.")
	}

	content, err := util.GetTemplate(CommitMessageTemplate, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(target, []byte(content), 0o755)
}

func Uninstall() error {
	g := git.New()
	hookPath, err := g.HookPath()
	if err != nil {
		return err
	}

	target := path.Join(strings.TrimSpace(hookPath), CommitMessageTemplate)
	if !file.IsFile(target) {
		return errors.New("hook file prepare-commit-msg is not exist.")
	}
	return os.Remove(target)
}
