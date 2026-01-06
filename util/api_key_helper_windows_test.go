//go:build windows

package util

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestCreateKillOnCloseJob(t *testing.T) {
	job, err := createKillOnCloseJob()
	if err != nil {
		t.Fatalf("createKillOnCloseJob() error = %v, want nil", err)
	}
	defer func() {
		_ = windows.CloseHandle(job)
	}()

	// Verify that the job handle is valid
	if job == 0 {
		t.Error("createKillOnCloseJob() returned invalid handle")
	}

	// Verify that KILL_ON_JOB_CLOSE flag is set
	var info windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	var returnLength uint32
	err = windows.QueryInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
		&returnLength,
	)
	if err != nil {
		t.Fatalf("QueryInformationJobObject() error = %v", err)
	}

	if info.BasicLimitInformation.LimitFlags&windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE == 0 {
		t.Error("KILL_ON_JOB_CLOSE flag is not set")
	}
}

func TestAssignProcessToJob_InvalidPID(t *testing.T) {
	job, err := createKillOnCloseJob()
	if err != nil {
		t.Fatalf("createKillOnCloseJob() error = %v", err)
	}
	defer func() {
		_ = windows.CloseHandle(job)
	}()

	tests := []struct {
		name string
		pid  int
	}{
		{
			name: "negative PID",
			pid:  -1,
		},
		{
			name: "PID exceeds max",
			pid:  0x80000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := assignProcessToJob(job, tt.pid)
			if err == nil {
				t.Error("assignProcessToJob() should return error for invalid PID")
			}
			if !strings.Contains(err.Error(), "invalid process ID") {
				t.Errorf("error should mention invalid PID, got: %v", err)
			}
		})
	}
}

func TestAssignProcessToJob_NonExistentPID(t *testing.T) {
	job, err := createKillOnCloseJob()
	if err != nil {
		t.Fatalf("createKillOnCloseJob() error = %v", err)
	}
	defer func() {
		_ = windows.CloseHandle(job)
	}()

	// Use a PID that likely doesn't exist (but is valid range)
	nonExistentPID := 99999

	_, err = assignProcessToJob(job, nonExistentPID)
	if err == nil {
		t.Error("assignProcessToJob() should return error for non-existent PID")
	}
}

func TestGetAPIKeyFromHelper_Windows_Success(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "simple echo command",
			command:  "echo test-api-key",
			expected: "test-api-key",
		},
		{
			name:     "command with whitespace",
			command:  "echo   test-key-with-spaces  ",
			expected: "test-key-with-spaces",
		},
		{
			name:     "powershell command",
			command:  `powershell -Command "Write-Output 'ps-key'"`,
			expected: "ps-key",
		},
		{
			name:     "set and echo variable",
			command:  "set KEY=win-key && echo %KEY%",
			expected: "win-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetAPIKeyFromHelper(context.Background(), tt.command)
			if err != nil {
				t.Fatalf("GetAPIKeyFromHelper() error = %v, want nil", err)
			}
			if result != tt.expected {
				t.Errorf("GetAPIKeyFromHelper() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetAPIKeyFromHelper_Windows_Timeout(t *testing.T) {
	// Use timeout command (Windows specific)
	// This will sleep for 15 seconds, which is longer than HelperTimeout (10s)
	command := "timeout /t 15 /nobreak >nul"

	start := time.Now()
	_, err := GetAPIKeyFromHelper(context.Background(), command)
	duration := time.Since(start)

	if err == nil {
		t.Fatal("GetAPIKeyFromHelper() should return timeout error")
	}

	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("error message should mention timeout, got: %v", err)
	}

	// Verify it actually timed out around the expected timeout duration
	// Allow up to 2 seconds margin
	if duration < HelperTimeout || duration > HelperTimeout+2*time.Second {
		t.Errorf("timeout duration = %v, want around %v", duration, HelperTimeout)
	}
}

func TestGetAPIKeyFromHelper_Windows_KillProcessTree(t *testing.T) {
	// Test that the Job Object kills the entire process tree
	// Create a command that spawns child processes
	command := `cmd /c "timeout /t 15 /nobreak >nul & timeout /t 15 /nobreak >nul"`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()
	_, err := GetAPIKeyFromHelper(ctx, command)
	duration := time.Since(start)

	if err == nil {
		t.Fatal("GetAPIKeyFromHelper() should return timeout error")
	}

	// Should timeout quickly (around 2 seconds, not 15)
	if duration > 3*time.Second {
		t.Errorf("timeout took too long: %v, expected around 2s", duration)
	}
}

func TestGetAPIKeyFromHelper_Windows_CommandFailure(t *testing.T) {
	tests := []struct {
		name    string
		command string
		wantErr string
	}{
		{
			name:    "non-existent command",
			command: "nonexistentcommand12345",
			wantErr: "failed",
		},
		{
			name:    "command with exit code 1",
			command: "exit 1",
			wantErr: "failed",
		},
		{
			name:    "invalid syntax",
			command: "echo %UNDEFINED_VAR && exit 1",
			wantErr: "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAPIKeyFromHelper(context.Background(), tt.command)
			if err == nil {
				t.Fatal("GetAPIKeyFromHelper() should return error for failed command")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error should contain %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetAPIKeyFromHelper_Windows_EmptyOutput(t *testing.T) {
	tests := []struct {
		name    string
		command string
	}{
		{
			name:    "command with no output",
			command: "rem no output",
		},
		{
			name:    "command outputting only whitespace",
			command: "echo    ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAPIKeyFromHelper(context.Background(), tt.command)
			if err == nil {
				t.Fatal("GetAPIKeyFromHelper() with empty output should return error")
			}
			if !strings.Contains(err.Error(), "empty output") {
				t.Errorf("error message should mention empty output, got: %v", err)
			}
		})
	}
}

func TestGetAPIKeyFromHelper_Windows_SecurityStderr(t *testing.T) {
	// Command that outputs to stderr (sensitive info should not be leaked in error)
	command := "echo secret-data 1>&2 && exit 1"

	_, err := GetAPIKeyFromHelper(context.Background(), command)
	if err == nil {
		t.Fatal("GetAPIKeyFromHelper() should return error when command fails")
	}

	// The error message should NOT contain the stderr output (security consideration)
	if strings.Contains(err.Error(), "secret-data") {
		t.Error("error message should not leak stderr content (security issue)")
	}
}

func TestGetAPIKeyFromHelper_Windows_ComplexCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "piped commands",
			command:  "echo my-api-key | findstr api",
			expected: "my-api-key",
		},
		{
			name:     "command with variable substitution",
			command:  "set KEY=test-123 && echo %KEY%",
			expected: "test-123",
		},
		{
			name:     "for loop",
			command:  `for /F %i in ('echo nested-key') do @echo %i`,
			expected: "nested-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetAPIKeyFromHelper(context.Background(), tt.command)
			if err != nil {
				t.Fatalf("GetAPIKeyFromHelper() error = %v, want nil", err)
			}
			if result != tt.expected {
				t.Errorf("GetAPIKeyFromHelper() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetAPIKeyFromHelper_Windows_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Start a long-running command
	done := make(chan error, 1)
	go func() {
		_, err := GetAPIKeyFromHelper(ctx, "timeout /t 30 /nobreak >nul")
		done <- err
	}()

	// Cancel after a short delay
	time.Sleep(500 * time.Millisecond)
	cancel()

	// Wait for the command to be cancelled
	select {
	case err := <-done:
		if err == nil {
			t.Error("GetAPIKeyFromHelper() should return error on context cancellation")
		}
		if !strings.Contains(err.Error(), "timeout") {
			t.Errorf("error should mention timeout, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("GetAPIKeyFromHelper() took too long to respond to cancellation")
	}
}

func TestGetAPIKeyFromHelper_Windows_MultipleInvocations(t *testing.T) {
	results := make(chan string, 3)
	errors := make(chan error, 3)

	for i := 0; i < 3; i++ {
		go func(n int) {
			result, err := GetAPIKeyFromHelper(
				context.Background(),
				fmt.Sprintf("echo test-key-%d", n),
			)
			if err != nil {
				errors <- err
			} else {
				results <- result
			}
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < 3; i++ {
		select {
		case result := <-results:
			if !strings.HasPrefix(result, "test-key-") {
				t.Errorf("unexpected result: %s", result)
			}
			successCount++
		case err := <-errors:
			t.Errorf("unexpected error: %v", err)
		case <-time.After(5 * time.Second):
			t.Error("timeout waiting for results")
		}
	}

	if successCount != 3 {
		t.Errorf("expected 3 successful invocations, got %d", successCount)
	}
}
