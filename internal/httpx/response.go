package httpx

import (
	"encoding/json"
	"net/http"
)

// WriteJSON encodes v as JSON and writes it to w with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return // 204 must not have a body
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Nothing more we can do once headers are written, but at least log it
		// if you have a logger accessible here — otherwise this is a safe no-op.
		_ = err
	}
}
