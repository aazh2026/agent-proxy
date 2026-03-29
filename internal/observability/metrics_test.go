package observability

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMetrics_RecordRequest(t *testing.T) {
	m := NewMetrics()

	m.RecordRequest(100*time.Millisecond, true)
	m.RecordRequest(200*time.Millisecond, true)
	m.RecordRequest(150*time.Millisecond, false)

	if m.GetTotalRequests() != 3 {
		t.Errorf("Expected 3 total requests, got %d", m.GetTotalRequests())
	}

	if m.GetSuccessRate() != 66.66666666666666 {
		t.Errorf("Expected 66.67%% success rate, got %f", m.GetSuccessRate())
	}

	if m.GetErrorRate() != 33.33333333333333 {
		t.Errorf("Expected 33.33%% error rate, got %f", m.GetErrorRate())
	}
}

func TestMetrics_GetAvgLatency(t *testing.T) {
	m := NewMetrics()

	m.RecordRequest(100*time.Millisecond, true)
	m.RecordRequest(200*time.Millisecond, true)
	m.RecordRequest(300*time.Millisecond, true)

	avg := m.GetAvgLatency()
	if avg <= 0 {
		t.Errorf("Expected positive avg latency, got %f", avg)
	}
}

func TestMetrics_GetQPS(t *testing.T) {
	m := NewMetrics()

	m.RecordRequest(100*time.Millisecond, true)
	m.RecordRequest(200*time.Millisecond, true)

	qps := m.GetQPS()
	if qps < 0 {
		t.Errorf("Expected non-negative QPS, got %f", qps)
	}
}

func TestMetrics_Reset(t *testing.T) {
	m := NewMetrics()

	m.RecordRequest(100*time.Millisecond, true)
	m.RecordRequest(200*time.Millisecond, false)

	m.Reset()

	if m.GetTotalRequests() != 0 {
		t.Errorf("Expected 0 total requests after reset, got %d", m.GetTotalRequests())
	}

	if m.GetSuccessRate() != 0 {
		t.Errorf("Expected 0 success rate after reset, got %f", m.GetSuccessRate())
	}
}

func TestMetrics_Snapshot(t *testing.T) {
	m := NewMetrics()

	m.RecordRequest(100*time.Millisecond, true)
	m.RecordRequest(200*time.Millisecond, false)

	snapshot := m.Snapshot()

	if snapshot.TotalRequests != 2 {
		t.Errorf("Expected 2 total requests in snapshot, got %d", snapshot.TotalRequests)
	}

	if snapshot.SuccessRate != 50.0 {
		t.Errorf("Expected 50%% success rate in snapshot, got %f", snapshot.SuccessRate)
	}

	if snapshot.ErrorRate != 50.0 {
		t.Errorf("Expected 50%% error rate in snapshot, got %f", snapshot.ErrorRate)
	}
}

func TestMetricsHandler_HandleMetrics(t *testing.T) {
	m := NewMetrics()
	handler := NewMetricsHandler(m, nil)

	m.RecordRequest(100*time.Millisecond, true)
	m.RecordRequest(200*time.Millisecond, false)

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.HandleMetrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Errorf("Expected text/plain content type, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "agent_proxy_requests_total") {
		t.Error("Expected agent_proxy_requests_total metric in response")
	}

	if !strings.Contains(body, "agent_proxy_qps") {
		t.Error("Expected agent_proxy_qps metric in response")
	}
}

func TestMetricsHandler_MethodNotAllowed(t *testing.T) {
	m := NewMetrics()
	handler := NewMetricsHandler(m, nil)

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/metrics", nil)
		w := httptest.NewRecorder()

		handler.HandleMetrics(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for %s, got %d", method, w.Code)
		}
	}
}
