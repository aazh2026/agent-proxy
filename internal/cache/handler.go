package cache

import (
	"encoding/json"
	"net/http"
)

type CacheHandler struct {
	cache *Cache
}

func NewCacheHandler(cache *Cache) *CacheHandler {
	return &CacheHandler{cache: cache}
}

func (h *CacheHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := h.cache.GetStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *CacheHandler) HandleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.cache.Clear()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *CacheHandler) HandleInvalidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	model := r.URL.Query().Get("model")
	if model == "" {
		http.Error(w, "model parameter required", http.StatusBadRequest)
		return
	}

	count := h.cache.InvalidateByPrefix(model)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "ok",
		"invalidated": count,
	})
}

func ShouldBypassCache(r *http.Request) bool {
	if r.Header.Get("X-Cache-Bypass") == "true" {
		return true
	}
	if r.Header.Get("Cache-Control") == "no-cache" {
		return true
	}
	return false
}
