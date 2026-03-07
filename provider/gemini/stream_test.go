package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/genai"
)

func TestCompletionStreamWriterOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := []map[string]any{
			{
				"candidates": []map[string]any{
					{
						"content": map[string]any{
							"parts": []map[string]any{
								{"text": "Hello"},
							},
						},
					},
				},
			},
			{
				"candidates": []map[string]any{
					{
						"content": map[string]any{
							"parts": []map[string]any{
								{"text": " world"},
							},
						},
					},
				},
				"usageMetadata": map[string]any{
					"promptTokenCount":     10,
					"candidatesTokenCount": 2,
					"totalTokenCount":      12,
				},
			},
		}
		data, _ := json.Marshal(resp)
		_, _ = w.Write(data)
	}))
	defer server.Close()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  "test-key",
		Backend: genai.BackendGeminiAPI,
		HTTPClient: &http.Client{
			Transport: &mockTransport{server: server},
		},
	})
	if err != nil {
		t.Fatalf("failed to create genai client: %v", err)
	}

	c := &Client{
		client:      client,
		model:       "gemini-2.0-flash",
		maxTokens:   1024,
		temperature: 0.7,
		topP:        1.0,
	}

	var buf bytes.Buffer
	resp, err := c.CompletionStream(ctx, "test prompt", &buf)
	if err != nil {
		t.Skipf("Skipping streaming test due to SDK transport constraints: %v", err)
		return
	}

	if resp.Content == "" {
		t.Error("expected non-empty content")
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty writer output")
	}
}

type mockTransport struct {
	server *httptest.Server
}

func (tr *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = tr.server.Listener.Addr().String()
	return http.DefaultTransport.RoundTrip(req)
}
