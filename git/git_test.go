package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCanExecuteGitDiff(t *testing.T) {
	// Test in the current directory (which is a git repository)
	t.Run("in git repository", func(t *testing.T) {
		cmd := New()
		ctx := context.Background()

		err := cmd.CanExecuteGitDiff(ctx)
		if err != nil {
			t.Errorf("CanExecuteGitDiff() should succeed in a git repository, got error: %v", err)
		}
	})

	// Test in a non-git directory
	t.Run("not in git repository", func(t *testing.T) {
		// Create a temporary directory
		tmpDir := t.TempDir()

		// Change to the temporary directory
		t.Chdir(tmpDir)

		cmd := New()
		ctx := context.Background()

		err := cmd.CanExecuteGitDiff(ctx)
		if err == nil {
			t.Error("CanExecuteGitDiff() should fail in a non-git directory")
		}

		expectedMsg := "not a git repository"
		if err != nil && err.Error() != expectedMsg {
			t.Logf("Got expected error: %v", err)
		}
	})

	// Test in a git repository subdirectory
	t.Run("in git repository subdirectory", func(t *testing.T) {
		// Create a test subdirectory in the git repository
		tmpDir := filepath.Join(".", "test_subdir")
		if err := os.MkdirAll(tmpDir, 0o755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Change to the subdirectory
		t.Chdir(tmpDir)

		cmd := New()
		ctx := context.Background()

		err := cmd.CanExecuteGitDiff(ctx)
		if err != nil {
			t.Errorf(
				"CanExecuteGitDiff() should succeed in a git repository subdirectory, got error: %v",
				err,
			)
		}
	})
}
