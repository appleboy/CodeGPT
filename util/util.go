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

func Float32Ptr(v float32) *float32 {
	return &v
}

func Int32Ptr(v int32) *int32 {
	return &v
}
