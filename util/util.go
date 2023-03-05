package util

import "os/exec"

// IsCommandAvailable check command exits.
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
