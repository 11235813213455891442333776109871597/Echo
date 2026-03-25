package echo

import (
	"encoding/json"
	"io"
	"net/http"
)

// Response is the structure echoed back to the caller.
type Response struct {
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers"`
	Query   map[string][]string `json:"query"`
	Body    string              `json:"body"`
}

// Handler echoes the incoming request's method, headers, query parameters, and body.
func Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	resp := Response{
		Method:  r.Method,
		Headers: r.Header,
		Query:   r.URL.Query(),
		Body:    string(bodyBytes),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
