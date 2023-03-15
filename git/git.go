package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
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
}

type Command struct {
	// Generate diffs with <n> lines of context instead of the usual three
	diffUnified int
	excludeList []string
}

func (c *Command) excludeFiles() []string {
	newFileLists := []string{}
	for _, f := range c.excludeList {
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
		"--unified=" + strconv.Itoa(c.diffUnified),
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
		"--signoff",
		fmt.Sprintf("--message=%s", val),
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

	target := path.Join(strings.TrimSpace(string(hookPath)), HookPrepareCommitMessageTemplate)
	if file.IsFile(target) {
		return errors.New("hook file prepare-commit-msg exist.")
	}

	content, err := util.GetTemplate(HookPrepareCommitMessageTemplate, nil)
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

	target := path.Join(strings.TrimSpace(string(hookPath)), HookPrepareCommitMessageTemplate)
	if !file.IsFile(target) {
		return errors.New("hook file prepare-commit-msg is not exist.")
	}
	return os.Remove(target)
}

func New(opts ...Option) *Command {
	cfg := &config{}

	// Loop through each option
	for _, o := range opts {
		// Call the option giving the instantiated
		o.apply(cfg)
	}

	return &Command{
		diffUnified: cfg.diffUnified,
		excludeList: append(excludeFromDiff, cfg.excludeList...),
	}
}
