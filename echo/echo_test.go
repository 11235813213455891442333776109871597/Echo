package echo_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/11235813213455891442333776109871597/Echo/echo"
)

func TestHandler_EchoesMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	rr := httptest.NewRecorder()

	echo.Handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp echo.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Method != http.MethodGet {
		t.Errorf("expected method GET, got %s", resp.Method)
	}
}

func TestHandler_EchoesBody(t *testing.T) {
	body := `{"hello":"world"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	echo.Handler(rr, req)

	var resp echo.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Body != body {
		t.Errorf("expected body %q, got %q", body, resp.Body)
	}
}

func TestHandler_EchoesQueryParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo?foo=bar&baz=qux", nil)
	rr := httptest.NewRecorder()

	echo.Handler(rr, req)

	var resp echo.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Query["foo"][0] != "bar" {
		t.Errorf("expected query foo=bar, got %v", resp.Query["foo"])
	}
}

func TestHandler_EchoesHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	req.Header.Set("X-Custom-Header", "test-value")
	rr := httptest.NewRecorder()

	echo.Handler(rr, req)

	var resp echo.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Headers["X-Custom-Header"]) == 0 || resp.Headers["X-Custom-Header"][0] != "test-value" {
		t.Errorf("expected X-Custom-Header=test-value, got %v", resp.Headers["X-Custom-Header"])
	}
}
