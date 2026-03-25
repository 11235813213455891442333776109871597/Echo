package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEchoHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?foo=bar&foo=baz", nil)
	req.Header.Set("X-Custom", "value")
	w := httptest.NewRecorder()

	echoHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	var resp EchoResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Method != http.MethodGet {
		t.Errorf("expected method GET, got %q", resp.Method)
	}
	if resp.Path != "/hello" {
		t.Errorf("expected path /hello, got %q", resp.Path)
	}
	if resp.Body != "" {
		t.Errorf("expected empty body, got %q", resp.Body)
	}
	if vals := resp.Query["foo"]; len(vals) != 2 || vals[0] != "bar" || vals[1] != "baz" {
		t.Errorf("unexpected query foo: %v", vals)
	}
	if vals := resp.Headers["X-Custom"]; len(vals) == 0 || vals[0] != "value" {
		t.Errorf("expected X-Custom header, got %v", resp.Headers["X-Custom"])
	}
}

func TestEchoHandler_POST_WithBody(t *testing.T) {
	body := `{"message":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	echoHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var resp EchoResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Method != http.MethodPost {
		t.Errorf("expected method POST, got %q", resp.Method)
	}
	if resp.Body != body {
		t.Errorf("expected body %q, got %q", body, resp.Body)
	}
}

func TestEchoHandler_NoQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	echoHandler(w, req)

	var resp EchoResponse
	body, _ := io.ReadAll(w.Result().Body)
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Query != nil {
		t.Errorf("expected nil query for request with no query params, got %v", resp.Query)
	}
}

func TestNewMux(t *testing.T) {
	mux := newMux()
	req := httptest.NewRequest(http.MethodGet, "/any/path", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 from mux, got %d", w.Code)
	}
}
