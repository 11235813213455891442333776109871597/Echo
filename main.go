package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/11235813213455891442333776109871597/Echo/auth"
	"github.com/11235813213455891442333776109871597/Echo/echo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Echo endpoints
	mux.HandleFunc("/echo", echo.Handler)

	// Auth endpoints
	mux.HandleFunc("/login", handleLogin)

	// Protected echo endpoint
	mux.Handle("/echo/secure", auth.Middleware(http.HandlerFunc(echo.Handler)))

	log.Printf("Echo server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := auth.Login(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
