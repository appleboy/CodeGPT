package util

import (
	"testing"

	"github.com/go-authgate/sdk-go/credstore"
)

// newTestCredStore returns a file-backed SecureStore using a temp directory.
// This avoids touching the OS keyring in CI/CD environments.
func newTestCredStore(t *testing.T) *credstore.SecureStore[string] {
	t.Helper()
	path := t.TempDir() + "/creds.json"
	file := credstore.NewStringFileStore(path)
	// Pass file as both primary and fallback so NewSecureStore always picks file.
	return credstore.NewSecureStore[string](file, file)
}

func TestCredStore_SetAndGet(t *testing.T) {
	store := newTestCredStore(t)

	if err := store.Save("openai.api_key", "sk-test-123"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	val, err := store.Load("openai.api_key")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if val != "sk-test-123" {
		t.Errorf("expected sk-test-123, got %s", val)
	}
}

func TestCredStore_GetMissing(t *testing.T) {
	store := newTestCredStore(t)

	_, err := store.Load("nonexistent.key")
	if err != credstore.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCredStore_Delete(t *testing.T) {
	store := newTestCredStore(t)

	if err := store.Save("gemini.api_key", "gm-test-456"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if err := store.Delete("gemini.api_key"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := store.Load("gemini.api_key")
	if err != credstore.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestGetCredential_Missing(t *testing.T) {
	// Save original and restore after test.
	original := credStore
	defer func() { credStore = original }()

	path := t.TempDir() + "/creds.json"
	file := credstore.NewStringFileStore(path)
	credStore = credstore.NewSecureStore[string](file, file)

	val, err := GetCredential("some.missing.key")
	if err != nil {
		t.Fatalf("expected nil error for missing key, got %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string for missing key, got %q", val)
	}
}

func TestSetAndGetCredential(t *testing.T) {
	// Save original and restore after test.
	original := credStore
	defer func() { credStore = original }()

	path := t.TempDir() + "/creds.json"
	file := credstore.NewStringFileStore(path)
	credStore = credstore.NewSecureStore[string](file, file)

	if err := SetCredential("openai.api_key", "sk-abc"); err != nil {
		t.Fatalf("SetCredential failed: %v", err)
	}

	val, err := GetCredential("openai.api_key")
	if err != nil {
		t.Fatalf("GetCredential failed: %v", err)
	}
	if val != "sk-abc" {
		t.Errorf("expected sk-abc, got %q", val)
	}
}

func TestDeleteCredential(t *testing.T) {
	// Save original and restore after test.
	original := credStore
	defer func() { credStore = original }()

	path := t.TempDir() + "/creds.json"
	file := credstore.NewStringFileStore(path)
	credStore = credstore.NewSecureStore[string](file, file)

	if err := SetCredential("gemini.api_key", "gm-xyz"); err != nil {
		t.Fatalf("SetCredential failed: %v", err)
	}

	if err := DeleteCredential("gemini.api_key"); err != nil {
		t.Fatalf("DeleteCredential failed: %v", err)
	}

	val, err := GetCredential("gemini.api_key")
	if err != nil {
		t.Fatalf("expected nil error after delete, got %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string after delete, got %q", val)
	}
}
