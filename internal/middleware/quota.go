package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/openclaw/agent-proxy/internal/auth"
	"github.com/openclaw/agent-proxy/internal/routing"
)

type QuotaMiddleware struct {
	tracker *routing.QuotaTracker
}

func NewQuotaMiddleware(tracker *routing.QuotaTracker) *QuotaMiddleware {
	return &QuotaMiddleware{
		tracker: tracker,
	}
}

func (m *QuotaMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := auth.GetUserID(r.Context())
		if userID == "" {
			userID = "default"
		}

		if err := m.tracker.CheckQuota(userID); err != nil {
			if err == routing.ErrQuotaExceeded {
				m.writeQuotaError(w, "Quota exceeded")
				return
			}
		}

		if err := m.tracker.IncrementRequests(userID, 1); err != nil {
			if err == routing.ErrQuotaExceeded {
				m.writeQuotaError(w, "Request quota exceeded")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *QuotaMiddleware) TrackTokens(userID string, tokens int64) error {
	return m.tracker.IncrementTokens(userID, tokens)
}

func (m *QuotaMiddleware) TrackCost(userID string, cost float64) error {
	return m.tracker.IncrementCost(userID, cost)
}

func (m *QuotaMiddleware) writeQuotaError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"type":    "quota_exceeded",
		},
	})
}

func (m *QuotaMiddleware) GetQuotaStatus(userID string) (map[string]interface{}, error) {
	counter := m.tracker.GetCounter(userID)

	status := map[string]interface{}{
		"user_id":  userID,
		"requests": counter.Requests,
		"tokens":   counter.Tokens,
		"cost":     counter.Cost,
		"reset_at": counter.ResetAt,
	}

	return status, nil
}
