//go:build !windows

package util

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// GetAPIKeyFromHelper executes a shell command to dynamically generate an API key.
// The command is executed in /bin/sh with a timeout controlled by the provided context.
// It returns the trimmed output from stdout, or an error if the command fails.
//
// On timeout, it kills the entire process group (shell and all descendants) using
// a two-phase approach: SIGTERM for graceful termination, then SIGKILL if needed.
//
// Security note: The returned API key is sensitive and should not be logged.
func GetAPIKeyFromHelper(ctx context.Context, helperCmd string) (string, error) {
	if helperCmd == "" {
		return "", fmt.Errorf("api_key_helper command is empty")
	}

	// Create context with timeout if not already set
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, HelperTimeout)
		defer cancel()
	}

	// Execute command in /bin/sh
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", helperCmd)

	// Create a new process group so we can kill all descendants on timeout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("api_key_helper start failed: %w", err)
	}

	// Wait for command completion in a goroutine
	done := make(chan error, 1)
	go func() {
		// Always Wait to avoid zombie processes
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		// Command completed normally
		if err != nil {
			// Don't include stderr in error message as it might contain sensitive info
			return "", fmt.Errorf("api_key_helper command failed: %w", err)
		}
		apiKey := strings.TrimSpace(stdout.String())
		if apiKey == "" {
			return "", fmt.Errorf("api_key_helper command returned empty output")
		}
		return apiKey, nil

	case <-ctx.Done():
		// Timeout or cancellation: terminate the process group gracefully, then forcefully
		if cmd.Process == nil {
			// Process handle not initialized; wait for cleanup and report timeout
			<-done
			return "", fmt.Errorf("api_key_helper command timeout after %v", HelperTimeout)
		}
		pgid := cmd.Process.Pid

		// First attempt: send SIGTERM to the entire process group for graceful shutdown
		_ = syscall.Kill(-pgid, syscall.SIGTERM)

		// Wait for graceful termination with a grace period
		select {
		case <-done:
			// Process exited after timeout was reached; treat as timeout regardless of exit status.
			// We intentionally ignore stdout/stderr here to avoid returning a key after a timeout.
			return "", fmt.Errorf("api_key_helper command timeout after %v", HelperTimeout)

		case <-time.After(2 * time.Second):
			// Grace period expired: send SIGKILL to force termination
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
			<-done // Wait for cleanup
			return "", fmt.Errorf("api_key_helper command timeout after %v", HelperTimeout)
		}
	}
}
