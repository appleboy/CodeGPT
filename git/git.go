package git

import (
	"errors"
	"log"
	"os/exec"
)

var excludeFromDiff = []string{
	"package-lock.json",
	"pnpm-lock.yaml",
	// yarn.lock, Cargo.lock, Gemfile.lock, Pipfile.lock, etc.
	"*.lock",
	"go.sum",
}

type Command struct {
	ignoreLists []string
}

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

// Diff compares the differences between two sets of data.
// It returns a string representing the differences and an error.
// If there are no differences, it returns an empty string and an error.
func (c *Command) DiffFiles() (string, error) {
	output, err := c.diffNames().Output()
	if err != nil {
		log.Fatal(err)
	}
	if string(output) == "" {
		return "", errors.New("please add your staged changes using git add <files...>")
	}

	output, err = c.diffFiles().Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(output), nil
}

// Hook to show git hook path
func (c *Command) HookPath() (string, error) {
	output, err := c.hookPath().Output()
	return string(output), err
}

func New() *Command {
	return &Command{}
}
