package util

import "testing"

func TestIsCommandAvailable(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "git command",
			args: args{
				cmd: "git",
			},
			want: true,
		},
		{
			name: "command not found",
			args: args{
				cmd: "foobar",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCommandAvailable(tt.args.cmd); got != tt.want {
				t.Errorf("IsCommandAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
