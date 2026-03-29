package observability

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type HealthChecker struct {
	db *sql.DB
}

func NewHealthChecker(db *sql.DB) *HealthChecker {
	return &HealthChecker{db: db}
}

func (h *HealthChecker) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
	})
}

func (h *HealthChecker) HandleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checks := make(map[string]interface{})
	allHealthy := true

	if err := h.checkDatabase(); err != nil {
		checks["database"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
		allHealthy = false
	} else {
		checks["database"] = map[string]interface{}{
			"status": "healthy",
		}
	}

	status := http.StatusOK
	if !allHealthy {
		status = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": map[bool]string{true: "ready", false: "not_ready"}[allHealthy],
		"checks": checks,
	})
}

func (h *HealthChecker) HandleLive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "alive",
	})
}

func (h *HealthChecker) checkDatabase() error {
	if h.db == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.db.PingContext(ctx)
}
