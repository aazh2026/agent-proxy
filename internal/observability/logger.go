package observability

import (
	"sync"
	"time"
)

type RequestLog struct {
	RequestID  string    `json:"request_id"`
	UserID     string    `json:"user_id"`
	Model      string    `json:"model"`
	Provider   string    `json:"provider"`
	StatusCode int       `json:"status_code"`
	LatencyMs  int64     `json:"latency_ms"`
	Error      string    `json:"error,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

type RequestLogger struct {
	mu      sync.RWMutex
	logs    []*RequestLog
	size    int
	maxSize int
}

func NewRequestLogger(maxSize int) *RequestLogger {
	return &RequestLogger{
		logs:    make([]*RequestLog, maxSize),
		maxSize: maxSize,
	}
}

func (l *RequestLogger) Log(entry *RequestLog) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs[l.size%l.maxSize] = entry
	l.size++
}

func (l *RequestLogger) GetLogs(limit int) []*RequestLog {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if limit <= 0 || limit > l.maxSize {
		limit = l.maxSize
	}

	start := 0
	if l.size > l.maxSize {
		start = l.size % l.maxSize
	}

	var result []*RequestLog
	for i := 0; i < limit && i < l.size; i++ {
		idx := (start + i) % l.maxSize
		if l.logs[idx] != nil {
			result = append(result, l.logs[idx])
		}
	}
	return result
}

func (l *RequestLogger) GetLogsByUser(userID string, limit int) []*RequestLog {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var result []*RequestLog
	for i := l.size - 1; i >= 0 && len(result) < limit; i-- {
		idx := i % l.maxSize
		if l.logs[idx] != nil && l.logs[idx].UserID == userID {
			result = append(result, l.logs[idx])
		}
	}
	return result
}

func (l *RequestLogger) GetLogsByModel(model string, limit int) []*RequestLog {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var result []*RequestLog
	for i := l.size - 1; i >= 0 && len(result) < limit; i-- {
		idx := i % l.maxSize
		if l.logs[idx] != nil && l.logs[idx].Model == model {
			result = append(result, l.logs[idx])
		}
	}
	return result
}

func (l *RequestLogger) GetLogsByStatus(statusCode int, limit int) []*RequestLog {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var result []*RequestLog
	for i := l.size - 1; i >= 0 && len(result) < limit; i-- {
		idx := i % l.maxSize
		if l.logs[idx] != nil && l.logs[idx].StatusCode == statusCode {
			result = append(result, l.logs[idx])
		}
	}
	return result
}

func (l *RequestLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = make([]*RequestLog, l.maxSize)
	l.size = 0
}
