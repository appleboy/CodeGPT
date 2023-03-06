package cmd

import (
	"errors"

	"github.com/appleboy/CodeGPT/hook"

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

		switch args[0] {
		case "install":
			if err := hook.Install(); err != nil {
				return err
			}

			color.Green("Install git hook: prepare-commit-msg successfully")
			color.Green("You can see the hook file: .git/hooks/prepare-commit-msg")
		case "uninstall":
			if err := hook.Uninstall(); err != nil {
				return err
			}

			color.Green("Remove git hook: prepare-commit-msg successfully")
		}

		return nil
	},
}
