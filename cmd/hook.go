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

// hookCmd represents the command for managing the prepare-commit-msg git hook.
var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage git prepare-commit-msg hook",
}

// hookInstallCmd installs the prepare-commit-msg hook.
var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install git prepare-commit-msg hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		g := git.New()

		if err := g.InstallHook(cmd.Context()); err != nil {
			return err
		}
		color.Green("Git hook 'prepare-commit-msg' installed successfully")
		color.Green("Hook file location: .git/hooks/prepare-commit-msg")

		return nil
	},
}

// hookUninstallCmd uninstalls the prepare-commit-msg hook.
var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall git prepare-commit-msg hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		g := git.New()

		if err := g.UninstallHook(cmd.Context()); err != nil {
			return err
		}
		color.Green("Git hook 'prepare-commit-msg' removed successfully")
		return nil
	},
}
