package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIsCommandAvailable(t *testing.T) {
	testCases := []struct {
		name  string
		cmd   string
		want  bool
		setup func() error
	}{
		{
			name: "command exists",
			cmd:  "ls",
			want: true,
		},
		{
			name: "command does not exist",
			cmd:  "nonexistentcommand",
			want: false,
		},
		{
			name: "command exists in path",
			cmd:  "git",
			want: true,
			setup: func() error {
				// Add /usr/local/bin to PATH for this test case
				return os.Setenv("PATH", "/usr/local/bin:"+os.Getenv("PATH"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				if err := tc.setup(); err != nil {
					t.Fatalf("failed to set up test case: %v", err)
				}
			}

			got := IsCommandAvailable(tc.cmd)

			if got != tc.want {
				t.Errorf("IsCommandAvailable(%q) = %v; want %v", tc.cmd, got, tc.want)
			}
		})
	}
}

func TestConvertToMap(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want Data
	}{
		{
			name: "convert slice to map",
			args: args{
				[]string{
					"TICKET_ID=ABC-1234",
					"Name=John Doe",
				},
			},
			want: Data{
				"Name":      "John Doe",
				"TICKET_ID": "ABC-1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToMap(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsGitRepo(t *testing.T) {
	type fields struct {
		setup func(t *testing.T)
		want  bool
	}

	testCases := []struct {
		name   string
		fields fields
	}{
		{
			name: "inside git repo root",
			fields: fields{
				want: true,
				setup: func(t *testing.T) {
					t.Helper()

					tmpDir := t.TempDir()

					cmd := exec.Command("git", "init")
					cmd.Dir = tmpDir
					if output, err := cmd.CombinedOutput(); err != nil {
						t.Fatalf("failed to init git repo: %v, output: %s", err, string(output))
					}

					if err := os.Chdir(tmpDir); err != nil {
						t.Fatalf("failed to chdir to repo root: %v", err)
					}
				},
			},
		},
		{
			name: "inside git repo subdir",
			fields: fields{
				want: true,
				setup: func(t *testing.T) {
					t.Helper()

					tmpDir := t.TempDir()

					cmd := exec.Command("git", "init")
					cmd.Dir = tmpDir
					if output, err := cmd.CombinedOutput(); err != nil {
						t.Fatalf("failed to init git repo: %v, output: %s", err, string(output))
					}

					subDir := filepath.Join(tmpDir, "subdir", "nested")
					if err := os.MkdirAll(subDir, 0o755); err != nil {
						t.Fatalf("failed to create subdir: %v", err)
					}

					if err := os.Chdir(subDir); err != nil {
						t.Fatalf("failed to chdir to subdir: %v", err)
					}
				},
			},
		},
		{
			name: "outside git repo",
			fields: fields{
				want: false,
				setup: func(t *testing.T) {
					t.Helper()

					tmpDir := t.TempDir()

					if err := os.Chdir(tmpDir); err != nil {
						t.Fatalf("failed to chdir to temp dir: %v", err)
					}
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			origWD, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get current working directory: %v", err)
			}
			t.Cleanup(func() {
				_ = os.Chdir(origWD)
			})

			if tc.fields.setup != nil {
				tc.fields.setup(t)
			}

			got := IsGitRepo()
			if got != tc.fields.want {
				t.Errorf("IsGitRepo() = %v; want %v", got, tc.fields.want)
			}
		})
	}
}
