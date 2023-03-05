package cmd

import (
	"fmt"
	"log"

	"github.com/appleboy/CodeGPT/git"
	"github.com/appleboy/CodeGPT/util"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Auto generate commit message",
	Run: func(cmd *cobra.Command, args []string) {
		// check git command exist
		if !util.IsCommandAvailable("git") {
			log.Fatal("To use CodeGPT, you must have git on your PATH")
		}
		diff, err := git.Diff()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(diff)
	},
}
