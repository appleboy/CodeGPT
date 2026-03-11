package util

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

const (
	// HelperTimeout is the maximum time to wait for the API key helper script to execute
	HelperTimeout = 10 * time.Second
	// DefaultRefreshInterval is the default interval for refreshing API keys
	DefaultRefreshInterval = 900 * time.Second // 15 minutes

	// helperKeyPrefix is the credstore key namespace for helper command cache entries.
	helperKeyPrefix = "helper:"
)

// apiKeyCache stores cached API keys with their metadata
type apiKeyCache struct {
	APIKey        string    `json:"apiKey"`
	LastFetchTime time.Time `json:"lastFetchTime"`
	HelperCmd     string    `json:"helperCmd"`
}

// helperCacheKey returns the credstore key for a given helper command.
func helperCacheKey(helperCmd string) string {
	hash := sha256.Sum256([]byte(helperCmd))
	return helperKeyPrefix + hex.EncodeToString(hash[:])
}

// readCache reads the cached API key from credstore.
func readCache(helperCmd string) (*apiKeyCache, error) {
	key := helperCacheKey(helperCmd)
	val, err := GetCredential(key)
	if err != nil {
		return nil, err
	}
	if val == "" {
		return nil, nil //nolint:nilnil // nil cache indicates cache miss, not an error
	}

	var cache apiKeyCache
	if err := json.Unmarshal([]byte(val), &cache); err != nil {
		return nil, err
	}

	// Verify the helper command matches
	if cache.HelperCmd != helperCmd {
		return nil, nil //nolint:nilnil // nil cache indicates cache miss, not an error
	}

	return &cache, nil
}

// writeCache writes the API key cache to credstore.
func writeCache(helperCmd, apiKey string) error {
	key := helperCacheKey(helperCmd)
	cache := apiKeyCache{
		APIKey:        apiKey,
		LastFetchTime: time.Now(),
		HelperCmd:     helperCmd,
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return SetCredential(key, string(data))
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
// Platform-specific implementations are in api_key_helper_unix.go and api_key_helper_windows.go.
//
// The command is executed with a timeout controlled by the provided context.
// It returns the trimmed output from stdout, or an error if the command fails.
//
// On timeout:
//   - Unix/Linux/macOS: kills the entire process group (shell and all descendants)
//   - Windows: terminates the Job Object (cmd.exe and all descendants)
//
// Security note: The returned API key is sensitive and should not be logged.

// GetAPIKeyFromHelperWithCache executes a shell command to dynamically generate an API key,
// with credstore-backed caching support. The API key is cached for the specified refresh interval.
// If refreshInterval is 0, the cache is disabled and the command is executed every time.
//
// The cache is stored in the OS keyring (macOS Keychain / Linux Secret Service /
// Windows Credential Manager) with a file-based fallback.
//
// Parameters:
//   - ctx: Context for controlling execution and timeouts
//   - helperCmd: The shell command to execute
//   - refreshInterval: How long to cache the API key (0 to disable caching)
//
// Returns the API key from cache if still valid, otherwise executes the helper command.
//
// Security note: The returned API key is sensitive and should not be logged.
func GetAPIKeyFromHelperWithCache(
	ctx context.Context,
	helperCmd string,
	refreshInterval time.Duration,
) (string, error) {
	if helperCmd == "" {
		return "", errors.New("api_key_helper command is empty")
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
