//go:build windows

package util

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// createKillOnCloseJob creates a Windows Job Object with KILL_ON_JOB_CLOSE flag.
// When the job handle is closed, all processes in the job will be terminated.
func createKillOnCloseJob() (windows.Handle, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}

	var info windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	// Enable KILL_ON_JOB_CLOSE flag
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	_, err = windows.SetInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)
	if err != nil {
		_ = windows.CloseHandle(job)
		return 0, err
	}
	return job, nil
}

// assignProcessToJob assigns a process to a Job Object.
// Returns the process handle which should be closed by the caller.
func assignProcessToJob(job windows.Handle, pid int) (windows.Handle, error) {
	// Validate PID range to prevent overflow
	if pid < 0 || pid > 0x7FFFFFFF {
		return 0, fmt.Errorf("invalid process ID: %d", pid)
	}

	// Get child process handle (requires PROCESS_ALL_ACCESS)
	// #nosec G115 -- PID validated above to prevent overflow
	hProc, err := windows.OpenProcess(windows.PROCESS_ALL_ACCESS, false, uint32(pid))
	if err != nil {
		return 0, err
	}
	// Assign to Job
	if err = windows.AssignProcessToJobObject(job, hProc); err != nil {
		_ = windows.CloseHandle(hProc)
		return 0, err
	}
	return hProc, nil
}

// GetAPIKeyFromHelper executes a shell command to dynamically generate an API key.
// The command is executed in cmd.exe with a timeout controlled by the provided context.
// It returns the trimmed output from stdout, or an error if the command fails.
//
// On timeout, it terminates the entire Job Object (cmd.exe and all descendants).
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

	// Execute command in cmd.exe
	cmd := exec.CommandContext(ctx, "cmd.exe", "/c", helperCmd)

	// Use CREATE_NEW_PROCESS_GROUP and CREATE_BREAKAWAY_FROM_JOB flags
	// This allows the child process to be assigned to a new Job,
	// even if the parent process is already in a Job
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP | windows.CREATE_BREAKAWAY_FROM_JOB,
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Create Job Object first
	job, err := createKillOnCloseJob()
	if err != nil {
		return "", fmt.Errorf("create job failed: %w", err)
	}
	// With KILL_ON_JOB_CLOSE, closing the job will kill all processes
	defer func() {
		_ = windows.CloseHandle(job)
	}()

	// Start the child process
	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("api_key_helper start failed: %w", err)
	}

	// Assign child process to Job
	hProc, err := assignProcessToJob(job, cmd.Process.Pid)
	if err != nil {
		// If unable to breakaway due to policy, fall back to just killing the process
		// (but this won't guarantee killing grandchild processes)
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return "", fmt.Errorf("assign process to job failed: %w", err)
	}
	defer func() {
		_ = windows.CloseHandle(hProc)
	}()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
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
		// Timeout: terminate the entire Job (all descendants)
		_ = windows.TerminateJobObject(job, 1)
		<-done // Wait for cleanup
		return "", fmt.Errorf("api_key_helper command timeout after %v", HelperTimeout)
	}
}
