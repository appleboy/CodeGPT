package cmd

import (
	"errors"

	"github.com/appleboy/CodeGPT/git"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "install/uninstall git prepare-commit-msg hook",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] != "install" && args[0] != "uninstall" {
			return errors.New("only support install or uninstall command")
		}

		g := git.New()

		switch args[0] {
		case "install":
			if err := g.InstallHook(); err != nil {
				return err
			}

			color.Green("Install git hook: prepare-commit-msg successfully")
			color.Green("You can see the hook file: .git/hooks/prepare-commit-msg")
		case "uninstall":
			if err := g.UninstallHook(); err != nil {
				return err
			}

			color.Green("Remove git hook: prepare-commit-msg successfully")
		}

		return nil
	},
}
