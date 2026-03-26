package observability

import (
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	mu                sync.RWMutex
	requests          uint64
	requestsSuccess   uint64
	requestsFailed    uint64
	totalLatency      uint64
	latencySamples    uint64
	requestsPerSecond map[string]uint64
	lastReset         time.Time
}

func NewMetrics() *Metrics {
	return &Metrics{
		requestsPerSecond: make(map[string]uint64),
		lastReset:         time.Now(),
	}
}

func (m *Metrics) RecordRequest(latency time.Duration, success bool) {
	atomic.AddUint64(&m.requests, 1)
	atomic.AddUint64(&m.totalLatency, uint64(latency.Milliseconds()))
	atomic.AddUint64(&m.latencySamples, 1)

	if success {
		atomic.AddUint64(&m.requestsSuccess, 1)
	} else {
		atomic.AddUint64(&m.requestsFailed, 1)
	}

	second := time.Now().Format("15:04:05")
	m.mu.Lock()
	m.requestsPerSecond[second]++
	m.mu.Unlock()
}

func (m *Metrics) GetQPS() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	elapsed := time.Since(m.lastReset).Seconds()
	if elapsed == 0 {
		return 0
	}

	total := uint64(0)
	for _, count := range m.requestsPerSecond {
		total += count
	}
	return float64(total) / elapsed
}

func (m *Metrics) GetAvgLatency() time.Duration {
	total := atomic.LoadUint64(&m.totalLatency)
	samples := atomic.LoadUint64(&m.latencySamples)
	if samples == 0 {
		return 0
	}
	return time.Duration(total/samples) * time.Millisecond
}

func (m *Metrics) GetSuccessRate() float64 {
	success := atomic.LoadUint64(&m.requestsSuccess)
	total := atomic.LoadUint64(&m.requests)
	if total == 0 {
		return 100.0
	}
	return float64(success) / float64(total) * 100.0
}

func (m *Metrics) GetErrorRate() float64 {
	failed := atomic.LoadUint64(&m.requestsFailed)
	total := atomic.LoadUint64(&m.requests)
	if total == 0 {
		return 0.0
	}
	return float64(failed) / float64(total) * 100.0
}

func (m *Metrics) GetTotalRequests() uint64 {
	return atomic.LoadUint64(&m.requests)
}

func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.StoreUint64(&m.requests, 0)
	atomic.StoreUint64(&m.requestsSuccess, 0)
	atomic.StoreUint64(&m.requestsFailed, 0)
	atomic.StoreUint64(&m.totalLatency, 0)
	atomic.StoreUint64(&m.latencySamples, 0)
	m.requestsPerSecond = make(map[string]uint64)
	m.lastReset = time.Now()
}

type MetricsSnapshot struct {
	TotalRequests uint64  `json:"total_requests"`
	SuccessRate   float64 `json:"success_rate"`
	ErrorRate     float64 `json:"error_rate"`
	AvgLatencyMs  int64   `json:"avg_latency_ms"`
	QPS           float64 `json:"qps"`
}

func (m *Metrics) Snapshot() *MetricsSnapshot {
	return &MetricsSnapshot{
		TotalRequests: m.GetTotalRequests(),
		SuccessRate:   m.GetSuccessRate(),
		ErrorRate:     m.GetErrorRate(),
		AvgLatencyMs:  m.GetAvgLatency().Milliseconds(),
		QPS:           m.GetQPS(),
	}
}
