package litellm

import (
	"testing"
	"time"
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
				WithToken("test-key"),
				WithModel("anthropic/claude-sonnet-4-6"),
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
				WithToken("test-key"),
				WithModel(""),
			),
			wantErr: errorsMissingModel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.valid(); err != tt.wantErr {
				t.Errorf("config.valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_config_defaults(t *testing.T) {
	cfg := newConfig()
	if cfg.baseURL != defaultBaseURL {
		t.Errorf("expected default baseURL %q, got %q", defaultBaseURL, cfg.baseURL)
	}
	if cfg.maxTokens != defaultMaxTokens {
		t.Errorf("expected default maxTokens %d, got %d", defaultMaxTokens, cfg.maxTokens)
	}
	if cfg.temperature != defaultTemperature {
		t.Errorf("expected default temperature %f, got %f", defaultTemperature, cfg.temperature)
	}
	if cfg.topP != defaultTopP {
		t.Errorf("expected default topP %f, got %f", defaultTopP, cfg.topP)
	}
}

func Test_config_options(t *testing.T) {
	cfg := newConfig(
		WithToken("sk-litellm-key"),
		WithModel("openai/gpt-4o"),
		WithBaseURL("http://myproxy:8000/v1"),
		WithProxyURL("http://proxy:8080"),
		WithSocksURL("socks5://proxy:1080"),
		WithTimeout(30*time.Second),
		WithMaxTokens(500),
		WithTemperature(0.7),
		WithTopP(0.9),
		WithFrequencyPenalty(0.5),
		WithPresencePenalty(0.3),
		WithSkipVerify(true),
		WithHeaders([]string{"X-Custom: value"}),
	)

	if cfg.token != "sk-litellm-key" {
		t.Errorf("expected token %q, got %q", "sk-litellm-key", cfg.token)
	}
	if cfg.model != "openai/gpt-4o" {
		t.Errorf("expected model %q, got %q", "openai/gpt-4o", cfg.model)
	}
	if cfg.baseURL != "http://myproxy:8000/v1" {
		t.Errorf("expected baseURL %q, got %q", "http://myproxy:8000/v1", cfg.baseURL)
	}
	if cfg.proxyURL != "http://proxy:8080" {
		t.Errorf("expected proxyURL %q, got %q", "http://proxy:8080", cfg.proxyURL)
	}
	if cfg.socksURL != "socks5://proxy:1080" {
		t.Errorf("expected socksURL %q, got %q", "socks5://proxy:1080", cfg.socksURL)
	}
	if cfg.timeout != 30*time.Second {
		t.Errorf("expected timeout %v, got %v", 30*time.Second, cfg.timeout)
	}
	if cfg.maxTokens != 500 {
		t.Errorf("expected maxTokens %d, got %d", 500, cfg.maxTokens)
	}
	if cfg.temperature != 0.7 {
		t.Errorf("expected temperature %f, got %f", float32(0.7), cfg.temperature)
	}
	if cfg.topP != 0.9 {
		t.Errorf("expected topP %f, got %f", float32(0.9), cfg.topP)
	}
	if cfg.frequencyPenalty != 0.5 {
		t.Errorf("expected frequencyPenalty %f, got %f", float32(0.5), cfg.frequencyPenalty)
	}
	if cfg.presencePenalty != 0.3 {
		t.Errorf("expected presencePenalty %f, got %f", float32(0.3), cfg.presencePenalty)
	}
	if !cfg.skipVerify {
		t.Error("expected skipVerify true")
	}
	if len(cfg.headers) != 1 || cfg.headers[0] != "X-Custom: value" {
		t.Errorf("expected headers [X-Custom: value], got %v", cfg.headers)
	}
}

func TestNew_missingToken(t *testing.T) {
	_, err := New(WithModel("test-model"))
	if err != errorsMissingToken {
		t.Errorf("expected errorsMissingToken, got %v", err)
	}
}

func TestNew_missingModel(t *testing.T) {
	_, err := New(WithToken("test-key"), WithModel(""))
	if err != errorsMissingModel {
		t.Errorf("expected errorsMissingModel, got %v", err)
	}
}

func TestNew_success(t *testing.T) {
	client, err := New(
		WithToken("test-key"),
		WithModel("anthropic/claude-sonnet-4-6"),
		WithBaseURL("http://localhost:4000/v1"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.model != "anthropic/claude-sonnet-4-6" {
		t.Errorf("expected model %q, got %q", "anthropic/claude-sonnet-4-6", client.model)
	}
}
