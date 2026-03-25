package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/11235813213455891442333776109871597/Echo/auth"
)

func TestLogin_ValidCredentials(t *testing.T) {
	token, err := auth.Login("admin", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	_, err := auth.Login("admin", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for invalid credentials")
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	_, err := auth.Login("unknown", "anything")
	if err == nil {
		t.Fatal("expected error for unknown user")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	token, err := auth.Login("user", "secret")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.Username != "user" {
		t.Errorf("expected username 'user', got %q", claims.Username)
	}
}

func TestValidateToken_Tampered(t *testing.T) {
	token, _ := auth.Login("admin", "password123")
	// Tamper with the payload segment
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatal("expected 3 JWT parts")
	}
	parts[1] = "dGFtcGVyZWQ" // base64("tampered")
	tampered := strings.Join(parts, ".")

	_, err := auth.ValidateToken(tampered)
	if err == nil {
		t.Fatal("expected error for tampered token")
	}
}

// TestValidateToken_NoneAlgorithm verifies that tokens using the "none"
// signing algorithm are rejected (algorithm-confusion / CVE fix).
func TestValidateToken_NoneAlgorithm(t *testing.T) {
	// A hand-crafted JWT signed with alg=none
	// Header: {"alg":"none","typ":"JWT"}
	// Payload: {"username":"admin","exp":9999999999}
	noneToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjo5OTk5OTk5OTk5fQ."

	_, err := auth.ValidateToken(noneToken)
	if err == nil {
		t.Fatal("expected error for none-algorithm token, but got nil (security bug!)")
	}
}

func TestMiddleware_MissingHeader(t *testing.T) {
	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestMiddleware_InvalidToken(t *testing.T) {
	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestMiddleware_ValidToken(t *testing.T) {
	token, _ := auth.Login("admin", "password123")

	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
