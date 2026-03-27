package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/openclaw/agent-proxy/internal/auth"
	"github.com/openclaw/agent-proxy/internal/logging"
)

type UsageTracker struct {
	db *sql.DB
}

type UsageRecord struct {
	UserID           string    `json:"user_id"`
	Model            string    `json:"model"`
	Provider         string    `json:"provider"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalTokens      int       `json:"total_tokens"`
	EstimatedCost    float64   `json:"estimated_cost"`
	Timestamp        time.Time `json:"timestamp"`
}

func NewUsageTracker(db *sql.DB) *UsageTracker {
	return &UsageTracker{db: db}
}

func (t *UsageTracker) RecordUsage(record *UsageRecord) error {
	_, err := t.db.Exec(
		`INSERT INTO usage_stats (user_id, model, provider, prompt_tokens, completion_tokens, total_tokens, estimated_cost, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		record.UserID, record.Model, record.Provider,
		record.PromptTokens, record.CompletionTokens, record.TotalTokens,
		record.EstimatedCost, record.Timestamp,
	)
	if err != nil {
		logging.Error("Failed to record usage: %v", err)
	}
	return err
}

func (t *UsageTracker) GetUserUsage(userID string, startTime, endTime time.Time) ([]*UsageRecord, error) {
	rows, err := t.db.Query(
		`SELECT user_id, model, provider, prompt_tokens, completion_tokens, total_tokens, estimated_cost, timestamp
		FROM usage_stats WHERE user_id = ? AND timestamp BETWEEN ? AND ? ORDER BY timestamp DESC`,
		userID, startTime, endTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*UsageRecord
	for rows.Next() {
		record := &UsageRecord{}
		err := rows.Scan(&record.UserID, &record.Model, &record.Provider,
			&record.PromptTokens, &record.CompletionTokens, &record.TotalTokens,
			&record.EstimatedCost, &record.Timestamp)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (t *UsageTracker) GetAggregatedUsage(userID string, startTime, endTime time.Time) (map[string]interface{}, error) {
	var totalRequests, totalPromptTokens, totalCompletionTokens, totalTokens int
	var totalCost float64

	err := t.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(prompt_tokens), 0), COALESCE(SUM(completion_tokens), 0),
		COALESCE(SUM(total_tokens), 0), COALESCE(SUM(estimated_cost), 0)
		FROM usage_stats WHERE user_id = ? AND timestamp BETWEEN ? AND ?`,
		userID, startTime, endTime,
	).Scan(&totalRequests, &totalPromptTokens, &totalCompletionTokens, &totalTokens, &totalCost)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_id":                 userID,
		"total_requests":          totalRequests,
		"total_prompt_tokens":     totalPromptTokens,
		"total_completion_tokens": totalCompletionTokens,
		"total_tokens":            totalTokens,
		"total_cost":              totalCost,
		"start_time":              startTime,
		"end_time":                endTime,
	}, nil
}

type UsageHandler struct {
	tracker *UsageTracker
}

func NewUsageHandler(tracker *UsageTracker) *UsageHandler {
	return &UsageHandler{tracker: tracker}
}

func (h *UsageHandler) HandleGetUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		userID = r.URL.Query().Get("user_id")
	}

	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	startTime := time.Now().AddDate(0, 0, -30)
	endTime := time.Now()

	if start := r.URL.Query().Get("start"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			startTime = t
		}
	}
	if end := r.URL.Query().Get("end"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			endTime = t
		}
	}

	aggregated, err := h.tracker.GetAggregatedUsage(userID, startTime, endTime)
	if err != nil {
		http.Error(w, "Failed to get usage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aggregated)
}
