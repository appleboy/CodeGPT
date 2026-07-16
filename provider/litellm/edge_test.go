package litellm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// --- Auth / API key errors ---

func TestCompletion_InvalidAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "Incorrect API key provided: sk-inva*****key.",
				"type":    "invalid_request_error",
				"code":    "invalid_api_key",
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("sk-invalid-key"),
		WithModel("anthropic/claude-sonnet-4-6"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Completion(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error for invalid API key, got nil")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected 401 in error, got: %v", err)
	}
}

func TestCompletionStream_InvalidAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "Invalid API key",
				"type":    "invalid_request_error",
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("sk-bad"),
		WithModel("openai/gpt-4o"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var buf bytes.Buffer
	_, err = client.CompletionStream(context.Background(), "test", &buf)
	if err == nil {
		t.Fatal("expected error for invalid API key on stream, got nil")
	}
}

// --- Model not found ---

func TestCompletion_ModelNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "Model 'nonexistent/model' not found. Check your LiteLLM proxy config.",
				"type":    "invalid_request_error",
				"code":    "model_not_found",
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("nonexistent/model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Completion(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error for model not found, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected 404 in error, got: %v", err)
	}
}

// --- Rate limit (429) ---

func TestCompletion_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "Rate limit reached for anthropic/claude-sonnet-4-6",
				"type":    "rate_limit_error",
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("anthropic/claude-sonnet-4-6"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Completion(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error for rate limit, got nil")
	}
	if !strings.Contains(err.Error(), "429") {
		t.Errorf("expected 429 in error, got: %v", err)
	}
}

// --- Context timeout ---

func TestCompletion_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": "late"}},
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("openai/gpt-4o"),
		WithBaseURL(server.URL+"/v1"),
		WithTimeout(100*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err = client.Completion(ctx, "test")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

// --- Empty / malformed responses ---

func TestCompletion_EmptyContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "chatcmpl-1",
			"object":  "chat.completion",
			"created": 1,
			"model":   "openai/gpt-4o",
			"choices": []map[string]any{
				{
					"index":         0,
					"message":       map[string]any{"role": "assistant", "content": ""},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]any{"prompt_tokens": 5, "completion_tokens": 0, "total_tokens": 5},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("openai/gpt-4o"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Completion(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "" {
		t.Errorf("expected empty content, got %q", resp.Content)
	}
}

func TestCompletion_NullContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(
			[]byte(
				`{"id":"1","object":"chat.completion","created":1,"model":"test","choices":[{"index":0,"message":{"role":"assistant","content":null},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":0,"total_tokens":1}}`,
			),
		)
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("test-model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Completion(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error on null content: %v", err)
	}
	if resp.Content != "" {
		t.Errorf("expected empty content for null, got %q", resp.Content)
	}
}

func TestCompletion_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("test-model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Completion(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

// --- Streaming edge cases ---

func TestCompletionStream_AbruptTermination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		fmt.Fprintf(
			w,
			"data: %s\n\n",
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"test","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
		)
		// Server closes connection without sending [DONE]
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("test-model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var buf bytes.Buffer
	resp, err := client.CompletionStream(context.Background(), "test", &buf)
	// Abrupt close triggers EOF which is handled gracefully
	if err != nil {
		t.Fatalf("unexpected error on abrupt stream: %v", err)
	}
	if resp.Content != "Hello" {
		t.Errorf("expected partial content %q, got %q", "Hello", resp.Content)
	}
}

func TestCompletionStream_EmptyChunks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		chunks := []string{
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"test","choices":[{"index":0,"delta":{"content":""},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"test","choices":[{"index":0,"delta":{"content":"World"},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"test","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":5,"completion_tokens":1,"total_tokens":6}}`,
		}
		for _, chunk := range chunks {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("test-model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var buf bytes.Buffer
	resp, err := client.CompletionStream(context.Background(), "test", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "World" {
		t.Errorf("expected %q, got %q", "World", resp.Content)
	}
}

func TestCompletionStream_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "Internal server error",
				"type":    "server_error",
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("test-model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var buf bytes.Buffer
	_, err = client.CompletionStream(context.Background(), "test", &buf)
	if err == nil {
		t.Fatal("expected error for server error on stream, got nil")
	}
}

// --- Reasoning content fallback ---

func TestCompletion_ReasoningContentFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(
			[]byte(
				`{"id":"1","object":"chat.completion","created":1,"model":"openai/o3","choices":[{"index":0,"message":{"role":"assistant","content":"","reasoning_content":"The answer is 4 because 2+2=4"},"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":8,"total_tokens":18}}`,
			),
		)
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("openai/o3"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Completion(context.Background(), "What is 2+2?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "The answer is 4 because 2+2=4" {
		t.Errorf("expected reasoning content fallback, got %q", resp.Content)
	}
}

func TestCompletionStream_ReasoningContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		chunks := []string{
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"openai/o3","choices":[{"index":0,"delta":{"reasoning_content":"thinking..."},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"openai/o3","choices":[{"index":0,"delta":{"content":"4"},"finish_reason":null}]}`,
			`{"id":"1","object":"chat.completion.chunk","created":1,"model":"openai/o3","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":5,"completion_tokens":2,"total_tokens":7}}`,
		}
		for _, chunk := range chunks {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("openai/o3"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var buf bytes.Buffer
	resp, err := client.CompletionStream(context.Background(), "2+2?", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "thinking...4" {
		t.Errorf("expected reasoning+content, got %q", resp.Content)
	}
}

// --- Model string format verification ---

func TestModelStringPassedCorrectly(t *testing.T) {
	var receivedModel string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)
		receivedModel, _ = req["model"].(string)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "1",
			"object":  "chat.completion",
			"created": 1,
			"model":   receivedModel,
			"choices": []map[string]any{
				{
					"index":         0,
					"message":       map[string]any{"role": "assistant", "content": "ok"},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]any{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer server.Close()

	models := []string{
		"anthropic/claude-sonnet-4-6",
		"openai/gpt-4o",
		"bedrock/anthropic.claude-sonnet-4-6-v1",
		"vertex_ai/gemini-2.5-flash",
		"groq/llama-4-scout-17b-16e-instruct",
	}

	for _, model := range models {
		t.Run(model, func(t *testing.T) {
			client, err := New(
				WithToken("test-key"),
				WithModel(model),
				WithBaseURL(server.URL+"/v1"),
			)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			_, err = client.Completion(context.Background(), "test")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if receivedModel != model {
				t.Errorf("model sent to proxy = %q, want %q", receivedModel, model)
			}
		})
	}
}

// --- API key forwarding ---

func TestAPIKeyForwarding(t *testing.T) {
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "1",
			"object":  "chat.completion",
			"created": 1,
			"model":   "test",
			"choices": []map[string]any{
				{
					"index":         0,
					"message":       map[string]any{"role": "assistant", "content": "ok"},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]any{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer server.Close()

	client, err := New(
		WithToken("sk-litellm-master-key-123"),
		WithModel("anthropic/claude-sonnet-4-6"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Completion(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedAuth != "Bearer sk-litellm-master-key-123" {
		t.Errorf("expected Bearer auth header, got %q", receivedAuth)
	}
}

// --- Context cancellation ---

func TestCompletion_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer server.Close()

	client, err := New(
		WithToken("test-key"),
		WithModel("test-model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.Completion(ctx, "test")
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

// --- GetSummaryPrefix with no tool calls in response ---

func TestGetSummaryPrefix_NoToolCallsReturnsContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "1",
			"object":  "chat.completion",
			"created": 1,
			"model":   "test",
			"choices": []map[string]any{
				{
					"index": 0,
					"message": map[string]any{
						"role":    "assistant",
						"content": "feat(auth): add OAuth2 support",
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
		WithToken("test-key"),
		WithModel("test-model"),
		WithBaseURL(server.URL+"/v1"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.GetSummaryPrefix(context.Background(), "test diff")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "feat(auth): add OAuth2 support" {
		t.Errorf("expected plain content fallback, got %q", resp.Content)
	}
}
