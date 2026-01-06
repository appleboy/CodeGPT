package util

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	// HelperTimeout is the maximum time to wait for the API key helper script to execute
	HelperTimeout = 10 * time.Second
	// DefaultRefreshInterval is the default interval for refreshing API keys
	DefaultRefreshInterval = 900 * time.Second // 15 minutes
)

// apiKeyCache stores cached API keys with their metadata
type apiKeyCache struct {
	APIKey        string    `json:"apiKey"`
	LastFetchTime time.Time `json:"lastFetchTime"`
	HelperCmd     string    `json:"helperCmd"`
}

// getCacheDir returns the cache directory path
func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	cacheDir := filepath.Join(home, ".config", "codegpt", ".cache")
	return cacheDir, nil
}

// getCacheFilePath returns the cache file path for a given helper command
func getCacheFilePath(helperCmd string) (string, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Use hash of helper command as filename
	hash := sha256.Sum256([]byte(helperCmd))
	filename := hex.EncodeToString(hash[:]) + ".json"
	return filepath.Join(cacheDir, filename), nil
}

// readCache reads the cached API key from file
func readCache(helperCmd string) (*apiKeyCache, error) {
	cachePath, err := getCacheFilePath(helperCmd)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Cache doesn't exist yet
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var cache apiKeyCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %w", err)
	}

	// Verify the helper command matches
	if cache.HelperCmd != helperCmd {
		return nil, nil // Cache is for a different command
	}

	return &cache, nil
}

// writeCache writes the API key cache to file
func writeCache(helperCmd, apiKey string) error {
	cachePath, err := getCacheFilePath(helperCmd)
	if err != nil {
		return err
	}

	cache := apiKeyCache{
		APIKey:        apiKey,
		LastFetchTime: time.Now(),
		HelperCmd:     helperCmd,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	// Write with restrictive permissions (only owner can read/write)
	if err := os.WriteFile(cachePath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// needsRefresh checks if the cached key needs to be refreshed
func needsRefresh(cache *apiKeyCache, refreshInterval time.Duration) bool {
	if cache == nil {
		return true
	}

	// Always refresh if interval is 0
	if refreshInterval == 0 {
		return true
	}

	// Check if cache is expired
	return time.Since(cache.LastFetchTime) >= refreshInterval
}

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
		pgid := cmd.Process.Pid

		// First attempt: send SIGTERM to the entire process group for graceful shutdown
		_ = syscall.Kill(-pgid, syscall.SIGTERM)

		// Wait for graceful termination with a grace period
		select {
		case err := <-done:
			if err != nil {
				return "", fmt.Errorf("api_key_helper terminated after timeout: %w", err)
			}
			apiKey := strings.TrimSpace(stdout.String())
			if apiKey == "" {
				return "", fmt.Errorf(
					"api_key_helper command returned empty output after timeout termination",
				)
			}
			return apiKey, nil

		case <-time.After(2 * time.Second):
			// Grace period expired: send SIGKILL to force termination
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
			<-done // Wait for cleanup
			return "", fmt.Errorf("api_key_helper command timed out after %v", HelperTimeout)
		}
	}
}

// GetAPIKeyFromHelperWithCache executes a shell command to dynamically generate an API key,
// with file-based caching support. The API key is cached for the specified refresh interval.
// If refreshInterval is 0, the cache is disabled and the command is executed every time.
//
// The cache is stored in ~/.config/codegpt/.cache/ directory with restrictive permissions (0600).
//
// Parameters:
//   - ctx: Context for controlling execution and timeouts
//   - helperCmd: The shell command to execute
//   - refreshInterval: How long to cache the API key (0 to disable caching)
//
// Returns the API key from cache if still valid, otherwise executes the helper command.
//
// Security note: The returned API key is sensitive and should not be logged.
// Cache files are stored with 0600 permissions but contain the API key in JSON format.
func GetAPIKeyFromHelperWithCache(
	ctx context.Context,
	helperCmd string,
	refreshInterval time.Duration,
) (string, error) {
	if helperCmd == "" {
		return "", fmt.Errorf("api_key_helper command is empty")
	}

	// Try to read from cache
	cache, err := readCache(helperCmd)
	if err != nil {
		// If cache read fails, log but continue to fetch fresh key
		// Don't fail the entire operation just because cache is broken
		cache = nil
	}

	// Check if we need to refresh
	if !needsRefresh(cache, refreshInterval) {
		return cache.APIKey, nil
	}

	// Fetch new API key
	apiKey, err := GetAPIKeyFromHelper(ctx, helperCmd)
	if err != nil {
		return "", err
	}

	// Write to cache (ignore errors to not block the operation)
	_ = writeCache(helperCmd, apiKey)

	return apiKey, nil
}
