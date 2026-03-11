package util

import (
	"os"
	"path/filepath"

	"github.com/go-authgate/sdk-go/credstore"
)

const credServiceName = "codegpt"

// credStore is the singleton SecureStore[string] instance.
// Initialized once; uses OS keyring with file-based fallback.
var credStore *credstore.SecureStore[string]

func init() {
	home, _ := os.UserHomeDir()
	fallbackPath := filepath.Join(home, ".config", "codegpt", ".cache", "credentials.json")
	keyring := credstore.NewStringKeyringStore(credServiceName)
	file := credstore.NewStringFileStore(fallbackPath)
	credStore = credstore.NewSecureStore(keyring, file)
}

// GetCredential retrieves a stored credential by key.
// Returns ("", nil) if not found.
func GetCredential(key string) (string, error) {
	val, err := credStore.Load(key)
	if err == credstore.ErrNotFound {
		return "", nil
	}
	return val, err
}

// SetCredential stores a credential by key.
func SetCredential(key, value string) error {
	return credStore.Save(key, value)
}

// DeleteCredential removes a credential by key.
func DeleteCredential(key string) error {
	return credStore.Delete(key)
}

// CredStoreIsKeyring reports whether the active backend is the OS keyring.
func CredStoreIsKeyring() bool {
	return credStore.UseKeyring()
}
