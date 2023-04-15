package util

import (
	"os"
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
