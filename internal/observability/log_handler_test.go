package observability

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/openclaw/agent-proxy/internal/auth"
)

func TestLogHandler_HandleGetLogs(t *testing.T) {
	logger := NewRequestLogger(10)
	handler := NewLogHandler(logger)

	entry := &RequestLog{
		RequestID:  "req1",
		UserID:     "user1",
		Model:      "model1",
		Provider:   "openai",
		StatusCode: 200,
		LatencyMs:  100,
		Timestamp:  time.Now(),
	}

	logger.Log(entry)

	req := httptest.NewRequest("GET", "/logs", nil)
	w := httptest.NewRecorder()

	handler.HandleGetLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var logs []*RequestLog
	if err := json.NewDecoder(w.Body).Decode(&logs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	if logs[0].RequestID != "req1" {
		t.Errorf("Expected first log req1, got %s", logs[0].RequestID)
	}
}

func TestLogHandler_HandleGetLogs_WithLimit(t *testing.T) {
	logger := NewRequestLogger(5)
	handler := NewLogHandler(logger)

	for i := 0; i < 5; i++ {
		entry := &RequestLog{
			RequestID:  fmt.Sprintf("req%d", i+1),
			UserID:     "user1",
			Model:      "model1",
			Provider:   "openai",
			StatusCode: 200,
			LatencyMs:  int64(100 * (i + 1)),
			Timestamp:  time.Now(),
		}
		logger.Log(entry)
	}

	req := httptest.NewRequest("GET", "/logs?limit=3", nil)
	w := httptest.NewRecorder()

	handler.HandleGetLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var logs []*RequestLog
	if err := json.NewDecoder(w.Body).Decode(&logs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(logs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(logs))
	}
}

func TestLogHandler_HandleGetLogs_ByUserID(t *testing.T) {
	logger := NewRequestLogger(10)
	handler := NewLogHandler(logger)

	entry1 := &RequestLog{
		RequestID:  "req1",
		UserID:     "user1",
		Model:      "model1",
		Provider:   "openai",
		StatusCode: 200,
		LatencyMs:  100,
		Timestamp:  time.Now(),
	}

	entry2 := &RequestLog{
		RequestID:  "req2",
		UserID:     "user2",
		Model:      "model2",
		Provider:   "anthropic",
		StatusCode: 200,
		LatencyMs:  150,
		Timestamp:  time.Now(),
	}

	logger.Log(entry1)
	logger.Log(entry2)

	req := httptest.NewRequest("GET", "/logs?user_id=user1", nil)
	w := httptest.NewRecorder()

	handler.HandleGetLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var logs []*RequestLog
	if err := json.NewDecoder(w.Body).Decode(&logs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log for user1, got %d", len(logs))
	}

	if logs[0].RequestID != "req1" {
		t.Errorf("Expected first log req1, got %s", logs[0].RequestID)
	}
}

func TestLogHandler_HandleClearLogs_Unauthorized(t *testing.T) {
	logger := NewRequestLogger(10)
	handler := NewLogHandler(logger)

	req := httptest.NewRequest("POST", "/logs/clear", nil)
	w := httptest.NewRecorder()

	handler.HandleClearLogs(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestLogHandler_HandleClearLogs_Authorized(t *testing.T) {
	logger := NewRequestLogger(10)
	handler := NewLogHandler(logger)

	entry := &RequestLog{
		RequestID:  "req1",
		UserID:     "user1",
		Model:      "model1",
		Provider:   "openai",
		StatusCode: 200,
		LatencyMs:  100,
		Timestamp:  time.Now(),
	}

	logger.Log(entry)

	req := httptest.NewRequest("POST", "/logs/clear", nil)
	req = req.WithContext(auth.SetUserID(req.Context(), "user1"))
	w := httptest.NewRecorder()

	handler.HandleClearLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}

	logs := logger.GetLogs(10)
	if len(logs) != 0 {
		t.Errorf("Expected 0 logs after clear, got %d", len(logs))
	}
}

func TestLogHandler_MethodNotAllowed(t *testing.T) {
	logger := NewRequestLogger(10)
	handler := NewLogHandler(logger)

	methods := []string{"PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/logs", nil)
		w := httptest.NewRecorder()

		handler.HandleGetLogs(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for %s, got %d", method, w.Code)
		}
	}

	req := httptest.NewRequest("GET", "/logs/clear", nil)
	w := httptest.NewRecorder()

	handler.HandleClearLogs(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for GET /logs/clear, got %d", w.Code)
	}
}
