package transport

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

// mockRoundTripper is used to intercept requests and record headers
type mockRoundTripper struct {
	lastReq *http.Request
	resp    *http.Response
	err     error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.lastReq = req
	if m.resp != nil || m.err != nil {
		return m.resp, m.err
	}
	return &http.Response{StatusCode: 200, Body: http.NoBody, Request: req}, nil
}

func TestDefaultHeaderTransport_CustomHeaders(t *testing.T) {
	mock := &mockRoundTripper{}
	tr := &DefaultHeaderTransport{
		Origin: mock,
		Header: http.Header{
			"X-Test": {"abc"},
			"Foo":    {"bar"},
		},
		AppName:    "myapp",
		AppVersion: "1.2.3",
	}
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	want := http.Header{
		"X-Test":        {"abc"},
		"Foo":           {"bar"},
		"X-App-Name":    {"myapp"},
		"X-App-Version": {"1.2.3"},
	}
	for k, v := range want {
		got := req.Header.Values(k)
		if !reflect.DeepEqual(got, v) {
			t.Errorf("Header %q = %v, want %v", k, got, v)
		}
	}
}

func TestDefaultHeaderTransport_EmptyHeadersAndAppInfo(t *testing.T) {
	mock := &mockRoundTripper{}
	tr := &DefaultHeaderTransport{
		Origin:     mock,
		Header:     http.Header{},
		AppName:    "",
		AppVersion: "",
	}
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	// x-app-name/version should not be present
	if req.Header.Get("x-app-name") != "" {
		t.Errorf("x-app-name should be empty")
	}
	if req.Header.Get("x-app-version") != "" {
		t.Errorf("x-app-version should be empty")
	}
}

func TestDefaultHeaderTransport_OriginErrorPropagation(t *testing.T) {
	mock := &mockRoundTripper{err: errors.New("mock error")}
	tr := &DefaultHeaderTransport{
		Origin:     mock,
		Header:     http.Header{},
		AppName:    "",
		AppVersion: "",
	}
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
	resp, err := tr.RoundTrip(req)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	if err == nil || err.Error() != "mock error" {
		t.Errorf("Expected error 'mock error', got %v", err)
	}
}

func TestDefaultHeaderTransport_MultipleHeaderValues(t *testing.T) {
	mock := &mockRoundTripper{}
	tr := &DefaultHeaderTransport{
		Origin: mock,
		Header: http.Header{
			"X-Multi": {"a", "b"},
		},
		AppName:    "app",
		AppVersion: "v",
	}
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	got := req.Header.Values("X-Multi")
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("X-Multi header = %v, want %v", got, want)
	}
}
