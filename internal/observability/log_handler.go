package observability

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/openclaw/agent-proxy/internal/auth"
)

type LogHandler struct {
	logger *RequestLogger
}

func NewLogHandler(logger *RequestLogger) *LogHandler {
	return &LogHandler{logger: logger}
}

func (h *LogHandler) HandleGetLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
	}

	userID := r.URL.Query().Get("user_id")
	model := r.URL.Query().Get("model")

	var logs []*RequestLog
	if userID != "" {
		logs = h.logger.GetLogsByUser(userID, limit)
	} else if model != "" {
		logs = h.logger.GetLogsByModel(model, limit)
	} else {
		logs = h.logger.GetLogs(limit)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *LogHandler) HandleClearLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.logger.Clear()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
