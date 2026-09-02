package openai

import (
	"errors"
	"testing"

	openai "github.com/sashabaranov/go-openai"

	"github.com/appleboy/CodeGPT/core"
)

func Test_config_valid(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config
		wantErr error
	}{
		{
			name: "valid config",
			cfg: newConfig(
				WithToken("test"),
				WithModel(openai.GPT3Dot5Turbo),
				WithProvider(core.OpenAI),
			),
			wantErr: nil,
		},
		{
			name:    "missing token",
			cfg:     newConfig(),
			wantErr: errorsMissingToken,
		},
		{
			name: "missing model",
			cfg: newConfig(
				WithToken("test"),
				WithModel(""),
				WithProvider(core.OpenAI),
			),
			wantErr: errorsMissingModel,
		},
		{
			name: "missing Azure deployment model",
			cfg: newConfig(
				WithToken("test"),
				WithModel(""),
				WithProvider(core.Azure),
			),
			wantErr: errorsMissingModel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg
			if err := cfg.valid(); !errors.Is(err, tt.wantErr) {
				t.Errorf("config.valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
