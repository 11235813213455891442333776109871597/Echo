package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// EchoResponse is the JSON structure returned by the echo handler.
type EchoResponse struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query,omitempty"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body,omitempty"`
}

// echoHandler reads the incoming request and reflects it back as JSON.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	resp := EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: map[string][]string(r.Header),
	}

	if len(r.URL.Query()) > 0 {
		resp.Query = map[string][]string(r.URL.Query())
	}

	if len(body) > 0 {
		resp.Body = string(body)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", echoHandler)
	return mux
}

func main() {
	addr := os.Getenv("ECHO_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Echo server listening on %s", addr)
	if err := http.ListenAndServe(addr, newMux()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
