package litellm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompletionStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		chunks := []string{
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"anthropic/claude-sonnet-4-6","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"anthropic/claude-sonnet-4-6","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"anthropic/claude-sonnet-4-6","choices":[{"index":0,"delta":{"content":" from"},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"anthropic/claude-sonnet-4-6","choices":[{"index":0,"delta":{"content":" LiteLLM"},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"anthropic/claude-sonnet-4-6","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":15,"completion_tokens":3,"total_tokens":18}}`,
		}

		for _, chunk := range chunks {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-token"),
		WithModel("anthropic/claude-sonnet-4-6"),
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

	expectedContent := "Hello from LiteLLM"
	if resp.Content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, resp.Content)
	}

	if buf.String() != expectedContent {
		t.Errorf("expected writer output %q, got %q", expectedContent, buf.String())
	}

	if resp.Usage.TotalTokens != 18 {
		t.Errorf("expected total tokens 18, got %d", resp.Usage.TotalTokens)
	}
	if resp.Usage.PromptTokens != 15 {
		t.Errorf("expected prompt tokens 15, got %d", resp.Usage.PromptTokens)
	}
	if resp.Usage.CompletionTokens != 3 {
		t.Errorf("expected completion tokens 3, got %d", resp.Usage.CompletionTokens)
	}
}

func TestCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{
			"id":      "chatcmpl-1",
			"object":  "chat.completion",
			"created": 1,
			"model":   "anthropic/claude-sonnet-4-6",
			"choices": []map[string]any{
				{
					"index": 0,
					"message": map[string]any{
						"role":    "assistant",
						"content": "4",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]any{
				"prompt_tokens":     10,
				"completion_tokens": 1,
				"total_tokens":      11,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-token"),
		WithModel("anthropic/claude-sonnet-4-6"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Completion(context.Background(), "What is 2+2?")
	if err != nil {
		t.Fatalf("Completion failed: %v", err)
	}

	if resp.Content != "4" {
		t.Errorf("expected content %q, got %q", "4", resp.Content)
	}

	if resp.Usage.TotalTokens != 11 {
		t.Errorf("expected total tokens 11, got %d", resp.Usage.TotalTokens)
	}
}

func TestGetSummaryPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{
			"id":      "chatcmpl-1",
			"object":  "chat.completion",
			"created": 1,
			"model":   "anthropic/claude-sonnet-4-6",
			"choices": []map[string]any{
				{
					"index": 0,
					"message": map[string]any{
						"role":    "assistant",
						"content": "",
						"tool_calls": []map[string]any{
							{
								"id":   "call_1",
								"type": "function",
								"function": map[string]any{
									"name":      "get_summary_prefix",
									"arguments": `{"prefix":"feat","scope":"litellm"}`,
								},
							},
						},
					},
					"finish_reason": "tool_calls",
				},
			},
			"usage": map[string]any{
				"prompt_tokens":     20,
				"completion_tokens": 5,
				"total_tokens":      25,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-token"),
		WithModel("anthropic/claude-sonnet-4-6"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.GetSummaryPrefix(context.Background(), "Added LiteLLM provider")
	if err != nil {
		t.Fatalf("GetSummaryPrefix failed: %v", err)
	}

	expected := "feat(litellm)"
	if resp.Content != expected {
		t.Errorf("expected content %q, got %q", expected, resp.Content)
	}

	if resp.Usage.TotalTokens != 25 {
		t.Errorf("expected total tokens 25, got %d", resp.Usage.TotalTokens)
	}
}

func TestGetSummaryPrefix_fallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode request to check if it has tools
		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)
		if _, hasTools := req["tools"]; hasTools {
			// Simulate function calling not supported
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"message": "tools not supported",
					"type":    "invalid_request_error",
				},
			})
			return
		}

		// Plain completion fallback
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "chatcmpl-1",
			"object":  "chat.completion",
			"created": 1,
			"model":   "openai/o1-mini",
			"choices": []map[string]any{
				{
					"index": 0,
					"message": map[string]any{
						"role":    "assistant",
						"content": "feat(litellm): add LiteLLM provider",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]any{
				"prompt_tokens":     10,
				"completion_tokens": 5,
				"total_tokens":      15,
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-token"),
		WithModel("openai/o1-mini"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.GetSummaryPrefix(context.Background(), "Added LiteLLM provider")
	if err != nil {
		t.Fatalf("GetSummaryPrefix fallback failed: %v", err)
	}

	if resp.Content != "feat(litellm): add LiteLLM provider" {
		t.Errorf("expected fallback content, got %q", resp.Content)
	}
}

func TestCompletion_noChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "chatcmpl-1",
			"object":  "chat.completion",
			"created": 1,
			"model":   "test-model",
			"choices": []map[string]any{},
			"usage": map[string]any{
				"prompt_tokens":     0,
				"completion_tokens": 0,
				"total_tokens":      0,
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-token"),
		WithModel("test-model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Completion(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error for empty choices")
	}
	if err.Error() != "no choices returned from LiteLLM proxy" {
		t.Errorf("expected 'no choices returned from LiteLLM proxy', got %q", err.Error())
	}
}
