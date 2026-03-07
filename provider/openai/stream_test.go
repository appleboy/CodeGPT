package openai

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompletionStream(t *testing.T) {
	// Create a mock SSE server that returns streaming chunks
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		chunks := []string{
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":" world"},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"gpt-4o","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":2,"total_tokens":12}}`,
		}

		for _, chunk := range chunks {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-token"),
		WithModel("gpt-4o"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var buf bytes.Buffer
	resp, err := client.CompletionStream(context.Background(), "test prompt", &buf)
	if err != nil {
		t.Fatalf("CompletionStream failed: %v", err)
	}

	expectedContent := "Hello world"
	if resp.Content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, resp.Content)
	}

	if buf.String() != expectedContent {
		t.Errorf("expected writer output %q, got %q", expectedContent, buf.String())
	}

	if resp.Usage.TotalTokens != 12 {
		t.Errorf("expected total tokens 12, got %d", resp.Usage.TotalTokens)
	}
}
