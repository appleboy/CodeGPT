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
	isAmend     bool
}

// excludeFiles returns a list of files to be excluded from git operations.
// It prepends each file in the excludeList with the exclude and top options.
func (c *Command) excludeFiles() []string {
	var excludedFiles []string
	for _, f := range c.excludeList {
		excludedFiles = append(excludedFiles, ":(exclude,top)"+f)
	}
	return excludedFiles
}

// diffNames generates the git command to list the names of changed files.
// It includes options to handle amended commits and staged changes.
func (c *Command) diffNames() *exec.Cmd {
	args := []string{
		"diff",
		"--name-only",
	}

	if c.isAmend {
		args = append(args, "HEAD^", "HEAD")
	} else {
		args = append(args, "--staged")
	}

	excludedFiles := c.excludeFiles()
	args = append(args, excludedFiles...)

	return exec.Command(
		"git",
		args...,
	)
}

// diffFiles generates the git command to show the differences between files.
// It includes options to ignore whitespace changes, use minimal diff algorithm,
// and set the number of context lines.
func (c *Command) diffFiles() *exec.Cmd {
	args := []string{
		"diff",
		"--ignore-all-space",
		"--diff-algorithm=minimal",
		"--unified=" + strconv.Itoa(c.diffUnified),
	}

	if c.isAmend {
		args = append(args, "HEAD^", "HEAD")
	} else {
		args = append(args, "--staged")
	}

	excludedFiles := c.excludeFiles()
	args = append(args, excludedFiles...)

	return exec.Command(
		"git",
		args...,
	)
}

// hookPath generates the git command to get the path of the hooks directory.
// This is used to locate where git hooks are stored.
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

// gitDir generates the git command to get the path of the git directory.
// This is used to determine the location of the .git directory.
func (c *Command) gitDir() *exec.Cmd {
	args := []string{
		"rev-parse",
		"--git-dir",
	}

	return exec.Command(
		"git",
		args...,
	)
}

// commit generates the git command to create a commit with the provided message.
// It includes options to skip pre-commit hooks, sign off the commit, and handle amendments.
func (c *Command) commit(val string) *exec.Cmd {
	args := []string{
		"commit",
		"--no-verify",
		"--signoff",
		fmt.Sprintf("--message=%s", val),
	}

	if c.isAmend {
		args = append(args, "--amend")
	}

	return exec.Command(
		"git",
		args...,
	)
}

// Commit creates a git commit with the provided message and returns the output or an error.
// It uses the commit method to generate the git command and execute it.
func (c *Command) Commit(val string) (string, error) {
	output, err := c.commit(val).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// GitDir to show the (by default, absolute) path of the git directory of the working tree.
func (c *Command) GitDir() (string, error) {
	output, err := c.gitDir().Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Diff compares the differences between two sets of data.
// It returns a string representing the differences and an error.
// If there are no differences, it returns an empty string and an error.
// DiffFiles compares the differences between two sets of data and returns the differences as a string and an error.
// It first lists the names of changed files and then shows the differences between them.
// If there are no staged changes, it returns an error message.
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

// InstallHook installs the prepare-commit-msg hook if it doesn't already exist.
// It retrieves the hooks directory path, checks if the hook file exists, and writes the hook file with executable permissions.
func (c *Command) InstallHook() error {
	hookPath, err := c.hookPath().Output()
	if err != nil {
		return err
	}

	target := path.Join(strings.TrimSpace(string(hookPath)), HookPrepareCommitMessageTemplate)
	if exists, err := file.IsFile(target); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if exists {
		return errors.New("hook file prepare-commit-msg exist")
	}

	content, err := util.GetTemplateByBytes(HookPrepareCommitMessageTemplate, nil)
	if err != nil {
		return err
	}

	// Write the hook file with executable permissions (0o755)
	return os.WriteFile(target, content, 0o755) //nolint:gosec
}

// UninstallHook removes the prepare-commit-msg hook if it exists.
// It retrieves the hooks directory path, checks if the hook file exists, and removes the hook file.
func (c *Command) UninstallHook() error {
	hookPath, err := c.hookPath().Output()
	if err != nil {
		return err
	}

	target := path.Join(strings.TrimSpace(string(hookPath)), HookPrepareCommitMessageTemplate)
	exists, err := file.IsFile(target)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("hook file prepare-commit-msg is not exist")
	}
	return os.Remove(target)
}

// New creates a new Command object with the provided options.
// It applies each option to the config object and initializes the Command object with the configurations.
func New(opts ...Option) *Command {
	// Instantiate a new config object with default values
	cfg := &config{}

	// Loop through each option passed as argument and apply it to the config object
	for _, o := range opts {
		o.apply(cfg)
	}

	// Instantiate a new Command object with the configurations from the config object
	cmd := &Command{
		diffUnified: cfg.diffUnified,
		// Append the user-defined excludeList to the default excludeFromDiff
		excludeList: append(excludeFromDiff, cfg.excludeList...),
		isAmend:     cfg.isAmend,
	}

	return cmd
}
