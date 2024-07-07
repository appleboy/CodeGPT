package cmd

import (
	"github.com/appleboy/CodeGPT/git"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
}

// hookCmd represents the command for installing/uninstalling the prepare-commit-msg hook.
var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "install/uninstall git prepare-commit-msg hook",
}

// hookInstallCmd installs the prepare-commit-msg hook.
var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install git prepare-commit-msg hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		g := git.New()

		if err := g.InstallHook(); err != nil {
			return err
		}
		color.Green("Install git hook: prepare-commit-msg successfully")
		color.Green("You can see the hook file: .git/hooks/prepare-commit-msg")

		return nil
	},
}

// hookUninstallCmd uninstalls the prepare-commit-msg hook.
var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall git prepare-commit-msg hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		g := git.New()

		if err := g.UninstallHook(); err != nil {
			return err
		}
		color.Green("Remove git hook: prepare-commit-msg successfully")
		return nil
	},
}
