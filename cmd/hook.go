package cmd

import (
	"errors"

	"github.com/appleboy/CodeGPT/hook"

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
			return hook.Install()
		case "uninstall":
			return hook.Uninstall()
		}

		return nil
	},
}
