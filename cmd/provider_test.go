package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/appleboy/CodeGPT/util"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

func TestNewLiteLLMCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	// Replace the OS keyring so tests never read or write real credentials.
	keyring.MockInit()
	zero, minute := 0, 60
	tests := []struct {
		name         string
		helper       string
		key          string
		sharedHelper string
		sharedKey    string
		want         string
		wantError    bool
		cached       bool
		interval     *int
	}{
		{
			name:         "default cache interval",
			sharedHelper: "echo fresh-default",
			cached:       true,
			want:         "cached-key",
		},
		{
			name:         "explicit cache interval",
			sharedHelper: "echo fresh-explicit",
			cached:       true,
			interval:     &minute,
			want:         "cached-key",
		},
		{
			name:         "disabled cache",
			sharedHelper: "echo fresh-disabled",
			cached:       true,
			interval:     &zero,
			want:         "fresh-disabled",
		},
		{name: "shared helper only", sharedHelper: "echo shared-helper", want: "shared-helper"},
		{
			name:         "shared helper beats static key",
			sharedHelper: "echo shared-helper",
			sharedKey:    "old-key",
			want:         "shared-helper",
		},
		{
			name:         "provider key beats shared helper",
			key:          "provider-key",
			sharedHelper: "exit 1",
			want:         "provider-key",
		},
		{
			name:         "provider helper beats all keys",
			helper:       "echo provider-helper",
			key:          "provider-key",
			sharedHelper: "exit 1",
			sharedKey:    "old-key",
			want:         "provider-helper",
		},
		{name: "shared static fallback", sharedKey: "shared-key", want: "shared-key"},
		{
			name:         "shared helper error",
			sharedHelper: "exit 1",
			sharedKey:    "old-key",
			wantError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			viper.Set("litellm.api_key_helper", tt.helper)
			viper.Set("litellm.api_key", tt.key)
			viper.Set("openai.api_key_helper", tt.sharedHelper)
			viper.Set("openai.api_key", tt.sharedKey)
			viper.Set("openai.model", "test-model")
			if tt.interval != nil {
				viper.Set("openai.api_key_helper_refresh_interval", *tt.interval)
			}
			if tt.cached {
				data, err := json.Marshal(
					map[string]any{"apiKey": "cached-key", "lastFetchTime": time.Now()},
				)
				if err != nil {
					t.Fatal(err)
				}
				cacheKey := fmt.Sprintf("helper:%x", sha256.Sum256([]byte(tt.sharedHelper)))
				if err := util.SetCredential(cacheKey, string(data)); err != nil {
					t.Fatal(err)
				}
			}
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if got := r.Header.Get("Authorization"); got != "Bearer "+tt.want {
						t.Errorf("Authorization = %q, want %q", got, "Bearer "+tt.want)
					}
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, `{"choices":[{"message":{"content":"ok"}}]}`)
				}),
			)
			defer server.Close()
			viper.Set("litellm.base_url", server.URL)
			client, err := NewLiteLLM(t.Context())
			if tt.wantError {
				if err == nil || !strings.Contains(err.Error(), "helper") {
					t.Fatalf("expected helper error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err := client.Completion(t.Context(), "test"); err != nil {
				t.Fatal(err)
			}
		})
	}
}
