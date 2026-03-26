package observability

import (
	"fmt"
	"net/http"
)

type MetricsHandler struct {
	metrics *Metrics
	logger  *RequestLogger
}

func NewMetricsHandler(metrics *Metrics, logger *RequestLogger) *MetricsHandler {
	return &MetricsHandler{
		metrics: metrics,
		logger:  logger,
	}
}

func (h *MetricsHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	snapshot := h.metrics.Snapshot()

	fmt.Fprintf(w, "# HELP agent_proxy_requests_total Total number of requests\n")
	fmt.Fprintf(w, "# TYPE agent_proxy_requests_total counter\n")
	fmt.Fprintf(w, "agent_proxy_requests_total %d\n", snapshot.TotalRequests)

	fmt.Fprintf(w, "# HELP agent_proxy_requests_success_rate Success rate percentage\n")
	fmt.Fprintf(w, "# TYPE agent_proxy_requests_success_rate gauge\n")
	fmt.Fprintf(w, "agent_proxy_requests_success_rate %.2f\n", snapshot.SuccessRate)

	fmt.Fprintf(w, "# HELP agent_proxy_requests_error_rate Error rate percentage\n")
	fmt.Fprintf(w, "# TYPE agent_proxy_requests_error_rate gauge\n")
	fmt.Fprintf(w, "agent_proxy_requests_error_rate %.2f\n", snapshot.ErrorRate)

	fmt.Fprintf(w, "# HELP agent_proxy_avg_latency_ms Average latency in milliseconds\n")
	fmt.Fprintf(w, "# TYPE agent_proxy_avg_latency_ms gauge\n")
	fmt.Fprintf(w, "agent_proxy_avg_latency_ms %d\n", snapshot.AvgLatencyMs)

	fmt.Fprintf(w, "# HELP agent_proxy_qps Queries per second\n")
	fmt.Fprintf(w, "# TYPE agent_proxy_qps gauge\n")
	fmt.Fprintf(w, "agent_proxy_qps %.2f\n", snapshot.QPS)
}

func (h *MetricsHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy"}`)
}
