package util

import (
	"os/exec"
	"strings"
)

// IsCommandAvailable checks whether a command is available in the PATH.
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ConvertToMap converts a slice of strings to a map.
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
