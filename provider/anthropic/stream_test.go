package anthropic

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liushuangls/go-anthropic/v2"
)

func TestCompletionStream(t *testing.T) {
	// Create a mock SSE server that returns Anthropic streaming events
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		events := []string{
			`event: message_start
data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","content":[],"model":"claude-sonnet-4-20250514","usage":{"input_tokens":10,"output_tokens":0}}}`,
			`event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
			`event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}`,
			`event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" world"}}`,
			`event: content_block_stop
data: {"type":"content_block_stop","index":0}`,
			`event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":2}}`,
			`event: message_stop
data: {"type":"message_stop"}`,
		}

		for _, event := range events {
			fmt.Fprintf(w, "%s\n\n", event)
		}
	}))
	defer server.Close()

	client := &Client{
		client: anthropic.NewClient(
			"test-token",
			anthropic.WithBaseURL(server.URL),
		),
		model:     anthropic.ModelClaude3Haiku20240307,
		maxTokens: 1024,
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

	if resp.Usage.PromptTokens != 10 {
		t.Errorf("expected prompt tokens 10, got %d", resp.Usage.PromptTokens)
	}

	if resp.Usage.TotalTokens != 12 {
		t.Errorf("expected total tokens 12, got %d", resp.Usage.TotalTokens)
	}
}
