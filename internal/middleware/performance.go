package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type PerformanceMiddleware struct {
}

func NewPerformanceMiddleware() *PerformanceMiddleware {
	return &PerformanceMiddleware{}
}

func (m *PerformanceMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		latency := time.Since(start)

		w.Header().Set("X-Response-Time", fmt.Sprintf("%dms", latency.Milliseconds()))

		w.Header().Set("X-Request-ID", r.Header.Get("X-Request-ID"))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
