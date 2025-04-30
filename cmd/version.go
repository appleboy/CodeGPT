package cmd

import (
	"fmt"

	"github.com/appleboy/CodeGPT/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the application version and commit information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version:", version.Version, "commit:", version.Commit)
	},
}
