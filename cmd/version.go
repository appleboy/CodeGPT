package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/appleboy/CodeGPT/version"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

/*
VersionInfo holds all version-related information for the application.
*/
type VersionInfo struct {
	Version   string `json:"version"`    // Application version
	GitCommit string `json:"git_commit"` // Git commit SHA
	BuildTime string `json:"build_time"` // Build timestamp
	GoVersion string `json:"go_version"` // Go language version
	BuildOS   string `json:"build_os"`   // Build operating system
	BuildArch string `json:"build_arch"` // Build architecture
	Platform  string `json:"platform"`   // Combined OS/Arch string
}

var outputFormat string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version, commit, build time, and environment details",
	RunE: func(cmd *cobra.Command, args []string) error {
		v := VersionInfo{
			Version:   version.Version,
			GitCommit: version.GitCommit,
			BuildTime: version.BuildTime,
			GoVersion: version.GoVersion,
			BuildOS:   version.BuildOS,
			BuildArch: version.BuildArch,
			Platform:  fmt.Sprintf("%s/%s", version.BuildOS, version.BuildArch),
		}
		return printVersion(outputFormat, v)
	},
}

func init() {
	versionCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text|json")
}

/*
shortCommit returns the first 7 characters of a Git commit SHA (short SHA).
*/
func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

/*
printVersion prints version information in the specified format.

format: "text" for colored CLI output, "json" for JSON output.
v:      VersionInfo struct containing version data.
*/
func printVersion(format string, v VersionInfo) error {
	// Use short SHA for Git commit
	shortV := v
	shortV.GitCommit = shortCommit(v.GitCommit)
	switch format {
	case "json":
		// Output as pretty-printed JSON
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(shortV)
	case "text":
		fallthrough
	default:
		// Output as colored CLI text
		blue := color.New(color.FgBlue, color.Bold)
		// Each row contains a label and its value
		rows := [][2]string{
			{"Version:", shortV.Version},
			{"Git Commit:", shortV.GitCommit},
			{"Build Time:", shortV.BuildTime},
			{"Go Version:", shortV.GoVersion},
			{"OS/Arch:", fmt.Sprintf("%s/%s", shortV.BuildOS, shortV.BuildArch)},
		}
		// Print each row with colored label and default value color
		for _, row := range rows {
			blue.Print(row[0])
			fmt.Printf(" %s\n", row[1])
		}
		return nil
	}
}
