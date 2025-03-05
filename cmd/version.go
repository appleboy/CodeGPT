package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string = ""
	Commit  string = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the application version and commit information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version:", Version, "commit:", Commit)
	},
}
