package api

import (
	"encoding/json"
	"net/http"
)

// ListExperiments returns a JSON array of configured AB experiments.
// This is a lightweight surface intended for MVP and to be wired to actual HTTP routing later.
func ListExperiments(w http.ResponseWriter, r *http.Request) {
	// Lightweight placeholder response for MVP
	type expResp struct {
		Name     string   `json:"name"`
		Variants []string `json:"variants"`
	}
	resp := struct {
		Experiments []expResp `json:"experiments"`
	}{Experiments: []expResp{}}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
