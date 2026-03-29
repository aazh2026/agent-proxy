package observability

import (
	"testing"
	"time"
)

func TestRequestLogger_LogAndGetLogs(t *testing.T) {
	logger := NewRequestLogger(5)

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

	logs := logger.GetLogs(10)
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}

	if logs[0].RequestID != "req1" {
		t.Errorf("Expected first log req1, got %s", logs[0].RequestID)
	}

	if logs[1].RequestID != "req2" {
		t.Errorf("Expected second log req2, got %s", logs[1].RequestID)
	}
}

func TestRequestLogger_GetLogsLimit(t *testing.T) {
	logger := NewRequestLogger(2)

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

	entry3 := &RequestLog{
		RequestID:  "req3",
		UserID:     "user3",
		Model:      "model3",
		Provider:   "google",
		StatusCode: 200,
		LatencyMs:  200,
		Timestamp:  time.Now(),
	}

	logger.Log(entry1)
	logger.Log(entry2)
	logger.Log(entry3)

	logs := logger.GetLogs(2)
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}

	// Due to circular buffer, we should get the last 2 logs added
	if logs[0].RequestID != "req2" {
		t.Errorf("Expected first log req2, got %s", logs[0].RequestID)
	}

	if logs[1].RequestID != "req3" {
		t.Errorf("Expected second log req3, got %s", logs[1].RequestID)
	}
}

func TestRequestLogger_GetLogsByUser(t *testing.T) {
	logger := NewRequestLogger(10)

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

	entry3 := &RequestLog{
		RequestID:  "req3",
		UserID:     "user1",
		Model:      "model3",
		Provider:   "google",
		StatusCode: 200,
		LatencyMs:  200,
		Timestamp:  time.Now(),
	}

	logger.Log(entry1)
	logger.Log(entry2)
	logger.Log(entry3)

	logs := logger.GetLogsByUser("user1", 10)
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs for user1, got %d", len(logs))
	}

	if logs[0].RequestID != "req1" {
		t.Errorf("Expected first log req1, got %s", logs[0].RequestID)
	}

	if logs[1].RequestID != "req3" {
		t.Errorf("Expected second log req3, got %s", logs[1].RequestID)
	}
}

func TestRequestLogger_GetLogsByModel(t *testing.T) {
	logger := NewRequestLogger(10)

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

	entry3 := &RequestLog{
		RequestID:  "req3",
		UserID:     "user1",
		Model:      "model1",
		Provider:   "google",
		StatusCode: 200,
		LatencyMs:  200,
		Timestamp:  time.Now(),
	}

	logger.Log(entry1)
	logger.Log(entry2)
	logger.Log(entry3)

	logs := logger.GetLogsByModel("model1", 10)
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs for model1, got %d", len(logs))
	}

	if logs[0].RequestID != "req1" {
		t.Errorf("Expected first log req1, got %s", logs[0].RequestID)
	}

	if logs[1].RequestID != "req3" {
		t.Errorf("Expected second log req3, got %s", logs[1].RequestID)
	}
}

func TestRequestLogger_GetLogsByStatus(t *testing.T) {
	logger := NewRequestLogger(10)

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
		StatusCode: 400,
		LatencyMs:  150,
		Timestamp:  time.Now(),
	}

	entry3 := &RequestLog{
		RequestID:  "req3",
		UserID:     "user1",
		Model:      "model3",
		Provider:   "google",
		StatusCode: 200,
		LatencyMs:  200,
		Timestamp:  time.Now(),
	}

	logger.Log(entry1)
	logger.Log(entry2)
	logger.Log(entry3)

	logs := logger.GetLogsByStatus(200, 10)
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs with status 200, got %d", len(logs))
	}

	if logs[0].RequestID != "req1" {
		t.Errorf("Expected first log req1, got %s", logs[0].RequestID)
	}

	if logs[1].RequestID != "req3" {
		t.Errorf("Expected second log req3, got %s", logs[1].RequestID)
	}
}

func TestRequestLogger_Clear(t *testing.T) {
	logger := NewRequestLogger(5)

	entry1 := &RequestLog{
		RequestID:  "req1",
		UserID:     "user1",
		Model:      "model1",
		Provider:   "openai",
		StatusCode: 200,
		LatencyMs:  100,
		Timestamp:  time.Now(),
	}

	logger.Log(entry1)

	logs := logger.GetLogs(10)
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	logger.Clear()

	logs = logger.GetLogs(10)
	if len(logs) != 0 {
		t.Errorf("Expected 0 logs after clear, got %d", len(logs))
	}
}
