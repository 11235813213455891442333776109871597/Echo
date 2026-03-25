package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret is used to sign and verify JWT tokens.
// It is loaded from the JWT_SECRET environment variable; a fallback value is
// used only when the variable is not set (development convenience only).
var jwtSecret = func() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte("change-me-in-production")
}()

// users is a simple in-memory credential store for demonstration purposes.
var users = map[string]string{
	"admin": "password123",
	"user":  "secret",
}

// Claims holds the JWT payload.
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login validates credentials and returns a signed JWT token.
func Login(username, password string) (string, error) {
	expected, ok := users[username]
	if !ok || expected != password {
		return "", errors.New("invalid credentials")
	}

	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken parses and validates a JWT token string.
// It explicitly requires HMAC signing to prevent algorithm-confusion attacks
// where an attacker supplies a token signed with the "none" algorithm or an
// asymmetric key whose public part is treated as the HMAC secret.
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			// Enforce HMAC signing — reject any other algorithm (e.g. "none", RSA).
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// Middleware is an HTTP middleware that requires a valid Bearer token.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		if _, err := ValidateToken(parts[1]); err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
