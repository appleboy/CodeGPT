package util

import (
	"bytes"
	"os/exec"
	"strings"
)

// IsCommandAvailable checks if a given command is available in the system's PATH.
// It takes a string argument 'cmd' which represents the command to check.
// It returns true if the command is found, otherwise false.
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ConvertToMap takes a slice of strings in the format "key=value" and converts it into a Data map.
// Each string in the slice is split into a key and value pair using the first occurrence of the "=" character.
// If a string does not contain the "=" character, it is ignored.
// The resulting map contains the keys and values from the input slice.
//
// Args:
//
//	args ([]string): A slice of strings where each string is in the format "key=value".
//
// Returns:
//
//	Data: A map where the keys and values are derived from the input slice.
func ConvertToMap(args []string) Data {
	m := make(Data)
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}
	return m
}

// IsGitRepo returns true if the current working directory is inside a Git work tree.
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false
	}

	return strings.TrimSpace(out.String()) == "true"
}
