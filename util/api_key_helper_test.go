package util

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetAPIKeyFromHelper_Success(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "simple echo command",
			command:  "echo 'test-api-key'",
			expected: "test-api-key",
		},
		{
			name:     "command with whitespace",
			command:  "echo '  test-key-with-spaces  '",
			expected: "test-key-with-spaces",
		},
		{
			name:     "command with newlines",
			command:  "printf 'key-with-newline\\n'",
			expected: "key-with-newline",
		},
		{
			name:     "multi-word echo",
			command:  "echo sk-1234567890abcdef",
			expected: "sk-1234567890abcdef",
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

func TestGetAPIKeyFromHelper_EmptyCommand(t *testing.T) {
	_, err := GetAPIKeyFromHelper(context.Background(), "")
	if err == nil {
		t.Fatal("GetAPIKeyFromHelper() with empty command should return error")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("error message should mention empty command, got: %v", err)
	}
}

func TestGetAPIKeyFromHelper_EmptyOutput(t *testing.T) {
	tests := []struct {
		name    string
		command string
	}{
		{
			name:    "true command with no output",
			command: "true",
		},
		{
			name:    "command outputting only whitespace",
			command: "echo '   '",
		},
		{
			name:    "command outputting only newlines",
			command: "printf '\\n\\n'",
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

func TestGetAPIKeyFromHelper_CommandFailure(t *testing.T) {
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
			name:    "command with syntax error",
			command: "if then fi",
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

func TestGetAPIKeyFromHelper_Timeout(t *testing.T) {
	// Command that sleeps longer than the timeout
	// The process group mechanism will kill both shell and sleep subprocess
	command := "sleep 15"

	start := time.Now()
	_, err := GetAPIKeyFromHelper(context.Background(), command)
	duration := time.Since(start)

	if err == nil {
		t.Fatal("GetAPIKeyFromHelper() should return timeout error")
	}
	// Error message can be either:
	// - "terminated after timeout" if SIGTERM succeeded
	// - "timed out" if SIGKILL was needed
	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "terminated") {
		t.Errorf("error message should mention timeout or termination, got: %v", err)
	}

	// Verify it actually timed out around the expected timeout duration
	// Allow up to 4 seconds margin (2s grace period + 2s buffer)
	if duration < HelperTimeout || duration > HelperTimeout+4*time.Second {
		t.Errorf("timeout duration = %v, want around %v", duration, HelperTimeout)
	}
}

func TestGetAPIKeyFromHelper_ComplexCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "piped commands",
			command:  "echo 'my-api-key' | tr '[:lower:]' '[:upper:]'",
			expected: "MY-API-KEY",
		},
		{
			name:     "command with variables",
			command:  "KEY=test-123; echo $KEY",
			expected: "test-123",
		},
		{
			name:     "command with subshell",
			command:  "echo $(printf 'nested-key')",
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

func TestGetAPIKeyFromHelper_SecurityStderr(t *testing.T) {
	// Command that outputs to stderr (sensitive info should not be leaked in error)
	command := "echo 'secret-data' >&2; exit 1"

	_, err := GetAPIKeyFromHelper(context.Background(), command)
	if err == nil {
		t.Fatal("GetAPIKeyFromHelper() should return error when command fails")
	}

	// The error message should NOT contain the stderr output (security consideration)
	if strings.Contains(err.Error(), "secret-data") {
		t.Error("error message should not leak stderr content (security issue)")
	}
}

func TestGetAPIKeyFromHelperWithCache_NoCaching(t *testing.T) {
	// Test with refreshInterval = 0 (no caching)
	command := "echo 'test-key-no-cache'"

	key1, err := GetAPIKeyFromHelperWithCache(context.Background(), command, 0)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}
	if key1 != "test-key-no-cache" {
		t.Errorf("Expected 'test-key-no-cache', got %q", key1)
	}

	// Second call should also execute (no caching)
	key2, err := GetAPIKeyFromHelperWithCache(context.Background(), command, 0)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}
	if key2 != "test-key-no-cache" {
		t.Errorf("Expected 'test-key-no-cache', got %q", key2)
	}
}

func TestGetAPIKeyFromHelperWithCache_WithCaching(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Use a counter file to generate different values each time the command runs
	counterFile := filepath.Join(tmpDir, "counter.txt")
	command := fmt.Sprintf(
		"f=%s; echo $(($(cat $f 2>/dev/null || echo 0) + 1)) | tee $f",
		counterFile,
	)

	// First call should execute and cache
	key1, err := GetAPIKeyFromHelperWithCache(context.Background(), command, 5*time.Second)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}

	// Small delay to ensure time difference
	time.Sleep(100 * time.Millisecond)

	// Second call should return cached value (same as first)
	key2, err := GetAPIKeyFromHelperWithCache(context.Background(), command, 5*time.Second)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}

	if key1 != key2 {
		t.Errorf("Cache should return same value: key1=%q, key2=%q", key1, key2)
	}
}

func TestGetAPIKeyFromHelperWithCache_CacheExpiration(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create a counter file that we'll update manually
	counterFile := filepath.Join(tmpDir, "counter2.txt")
	command := fmt.Sprintf("cat %s", counterFile)

	// Write initial value
	if err := os.WriteFile(counterFile, []byte("value1"), 0o600); err != nil {
		t.Fatalf("Failed to write counter file: %v", err)
	}

	// First call with short refresh interval
	key1, err := GetAPIKeyFromHelperWithCache(context.Background(), command, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}

	// Update the file with a different value
	if err := os.WriteFile(counterFile, []byte("value2"), 0o600); err != nil {
		t.Fatalf("Failed to update counter file: %v", err)
	}

	// Wait for cache to expire
	time.Sleep(600 * time.Millisecond)

	// Second call should fetch fresh value
	key2, err := GetAPIKeyFromHelperWithCache(context.Background(), command, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}

	// Keys should be different (cache expired)
	if key1 == key2 {
		t.Errorf("Cache should have expired and returned new value: key1=%q, key2=%q", key1, key2)
	}
	if key1 != "value1" {
		t.Errorf("First key should be 'value1', got %q", key1)
	}
	if key2 != "value2" {
		t.Errorf("Second key should be 'value2', got %q", key2)
	}
}

func TestGetAPIKeyFromHelperWithCache_DifferentCommands(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	cmd1 := "echo 'key-one'"
	cmd2 := "echo 'key-two'"

	// Get keys from different commands
	key1, err := GetAPIKeyFromHelperWithCache(context.Background(), cmd1, 5*time.Second)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}

	key2, err := GetAPIKeyFromHelperWithCache(context.Background(), cmd2, 5*time.Second)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}

	// Keys should be different (different commands)
	if key1 == key2 {
		t.Errorf("Different commands should return different keys: key1=%q, key2=%q", key1, key2)
	}

	if key1 != "key-one" {
		t.Errorf("Expected 'key-one', got %q", key1)
	}
	if key2 != "key-two" {
		t.Errorf("Expected 'key-two', got %q", key2)
	}
}

func TestGetAPIKeyFromHelperWithCache_CacheFilePermissions(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	command := "echo 'test-permissions'"

	// Execute command to create cache file
	_, err := GetAPIKeyFromHelperWithCache(context.Background(), command, 5*time.Second)
	if err != nil {
		t.Fatalf("GetAPIKeyFromHelperWithCache() error = %v", err)
	}

	// Check cache file permissions
	cachePath, err := getCacheFilePath(command)
	if err != nil {
		t.Fatalf("getCacheFilePath() error = %v", err)
	}

	info, err := os.Stat(cachePath)
	if err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}

	// Check that file has restrictive permissions (0600)
	perm := info.Mode().Perm()
	if perm != 0o600 {
		t.Errorf("Cache file should have 0600 permissions, got %o", perm)
	}
}
