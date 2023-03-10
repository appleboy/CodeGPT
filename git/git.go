package git

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/appleboy/CodeGPT/util"
	"github.com/appleboy/com/file"
)

var excludeFromDiff = []string{
	"package-lock.json",
	"pnpm-lock.yaml",
	// yarn.lock, Cargo.lock, Gemfile.lock, Pipfile.lock, etc.
	"*.lock",
	"go.sum",
	"go.mod",
}

type Command struct{}

func (c *Command) excludeFiles() []string {
	newFileLists := []string{}
	for _, f := range excludeFromDiff {
		newFileLists = append(newFileLists, ":(exclude)"+f)
	}

	return newFileLists
}

func (c *Command) diffNames() *exec.Cmd {
	args := []string{
		"diff",
		"--staged",
		"--name-only",
	}

	args = append(args, c.excludeFiles()...)

	return exec.Command(
		"git",
		args...,
	)
}

func (c *Command) diffFiles() *exec.Cmd {
	args := []string{
		"diff",
		"--staged",
		"--ignore-all-space",
		"--diff-algorithm=minimal",
		"--function-context",
	}

	args = append(args, c.excludeFiles()...)

	return exec.Command(
		"git",
		args...,
	)
}

func (c *Command) hookPath() *exec.Cmd {
	args := []string{
		"rev-parse",
		"--git-path",
		"hooks",
	}

	return exec.Command(
		"git",
		args...,
	)
}

func (c *Command) commit(val string) *exec.Cmd {
	args := []string{
		"commit",
		"--no-verify",
		"--file",
		val,
	}

	return exec.Command(
		"git",
		args...,
	)
}

func (c *Command) Commit(val string) (string, error) {
	output, err := c.commit(val).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Diff compares the differences between two sets of data.
// It returns a string representing the differences and an error.
// If there are no differences, it returns an empty string and an error.
func (c *Command) DiffFiles() (string, error) {
	output, err := c.diffNames().Output()
	if err != nil {
		return "", err
	}
	if string(output) == "" {
		return "", errors.New("please add your staged changes using git add <files...>")
	}

	output, err = c.diffFiles().Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (c *Command) InstallHook() error {
	hookPath, err := c.hookPath().Output()
	if err != nil {
		return err
	}

	target := path.Join(strings.TrimSpace(string(hookPath)), CommitMessageTemplate)
	if file.IsFile(target) {
		return errors.New("hook file prepare-commit-msg exist.")
	}

	content, err := util.GetTemplate(CommitMessageTemplate, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(target, []byte(content), 0o755)
}

func (c *Command) UninstallHook() error {
	hookPath, err := c.hookPath().Output()
	if err != nil {
		return err
	}

	target := path.Join(strings.TrimSpace(string(hookPath)), CommitMessageTemplate)
	if !file.IsFile(target) {
		return errors.New("hook file prepare-commit-msg is not exist.")
	}
	return os.Remove(target)
}

func New() *Command {
	return &Command{}
}
