package openai

import (
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func Test_config_vaild(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config
		wantErr error
	}{
		{
			name:    "test vaild",
			cfg:     newConfig(),
			wantErr: errorsMissingToken,
		},
		{
			name: "model vaild",
			cfg: newConfig(
				WithToken("test"),
				WithModel("test"),
			),
			wantErr: errorsMissingModel,
		},
		{
			name: "missing Azure deployment model",
			cfg: newConfig(
				WithToken("test"),
				WithModel(openai.GPT3Dot5Turbo),
				WithProvider(AZURE),
			),
			wantErr: errorsMissingAzureModel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg
			if err := cfg.vaild(); err != tt.wantErr {
				t.Errorf("config.vaild() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
