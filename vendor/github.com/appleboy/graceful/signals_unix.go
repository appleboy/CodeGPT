//go:build !windows
// +build !windows

package graceful

import (
	"os"
	"syscall"
)

var signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP}
