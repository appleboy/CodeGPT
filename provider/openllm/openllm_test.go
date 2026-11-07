package openllm

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew_missingModel(t *testing.T) {
	if _, err := New(WithToken("k")); err == nil {
		t.Fatal("expected error for missing model")
	}
}

func TestNew_keylessLocalDefaults(t *testing.T) {
	c, err := New(WithModel("lite"))
	if err != nil {
		t.Fatalf("keyless client must be valid for the local daemon: %v", err)
	}
	if c.model != "lite" {
		t.Fatalf("model = %q", c.model)
	}
	cfg := newConfig()
	if cfg.baseURL != "http://127.0.0.1:8787/v1" {
		t.Fatalf("default base URL = %q", cfg.baseURL)
	}
}

func TestWithBaseURL_emptyKeepsDefault(t *testing.T) {
	cfg := newConfig(WithBaseURL(""))
	if cfg.baseURL != defaultBaseURL {
		t.Fatalf("empty base URL must keep the default, got %q", cfg.baseURL)
	}
	cfg = newConfig(WithBaseURL("https://www.openllm.sh/v1"))
	if cfg.baseURL != "https://www.openllm.sh/v1" {
		t.Fatalf("base URL = %q", cfg.baseURL)
	}
}

func newServer(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := New(
		WithToken("sk-llm-test.secret"),
		WithModel("lite"),
		WithBaseURL(srv.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestCompletion(t *testing.T) {
	c := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer sk-llm-test.secret" {
			t.Errorf("Authorization header mismatch")
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(
			w,
			`{"choices":[{"message":{"content":"feat: hello"}}],"usage":{"total_tokens":7}}`,
		)
	})
	resp, err := c.Completion(t.Context(), "diff")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "feat: hello" || resp.Usage.TotalTokens != 7 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestCompletion_noChoices(t *testing.T) {
	c := newServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"choices":[]}`)
	})
	if _, err := c.Completion(t.Context(), "diff"); err == nil {
		t.Fatal("expected error when no choices are returned")
	}
}

func TestCompletion_authError(t *testing.T) {
	c := newServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(
			w,
			`{"error":{"message":"[openllm] API key required. Run openllm login.","type":"auth_error"}}`,
		)
	})
	_, err := c.Completion(t.Context(), "diff")
	if err == nil || !strings.Contains(err.Error(), "openllm login") {
		t.Fatalf("expected the daemon auth hint to propagate, got %v", err)
	}
}

func TestCompletionStream(t *testing.T) {
	c := newServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"fix\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\": ok\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[],\"usage\":{\"total_tokens\":5}}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	})
	var buf bytes.Buffer
	resp, err := c.CompletionStream(t.Context(), "diff", &buf)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != "fix: ok" || resp.Content != "fix: ok" || resp.Usage.TotalTokens != 5 {
		t.Fatalf("stream = %q, resp = %+v", buf.String(), resp)
	}
}

func TestGetSummaryPrefix(t *testing.T) {
	c := newServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(
			w,
			`{"choices":[{"message":{"tool_calls":[{"id":"1","type":"function","function":{"name":"get_summary_prefix","arguments":"{\"prefix\":\"feat\",\"scope\":\"cmd\"}"}}]}}]}`,
		)
	})
	resp, err := c.GetSummaryPrefix(t.Context(), "diff")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "feat(cmd)" {
		t.Fatalf("content = %q", resp.Content)
	}
}

func TestGetSummaryPrefix_fallback(t *testing.T) {
	calls := 0
	c := newServer(t, func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls == 1 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error":{"message":"tools unsupported"}}`)
			return
		}
		fmt.Fprint(w, `{"choices":[{"message":{"content":"chore"}}]}`)
	})
	resp, err := c.GetSummaryPrefix(t.Context(), "diff")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "chore" || calls != 2 {
		t.Fatalf("content = %q, calls = %d", resp.Content, calls)
	}
}
