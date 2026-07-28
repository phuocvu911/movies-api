package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Single, shared validator instance — create once, reuse everywhere.
// It's safe for concurrent use.

// writeJSON serializes v to the response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("write response: %v", err)
		}
	}
}

// decodeJSON decode json body into struct variable
func decodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}
