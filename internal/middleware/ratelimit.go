package middleware

import (
	"fmt"
	"net/http"

	"github.com/openclaw/agent-proxy/internal/routing"
)

type RateLimitMiddleware struct {
	userLimiter   *routing.RateLimiter
	ipLimiter     *routing.RateLimiter
	globalLimiter *routing.RateLimiter
}

func NewRateLimitMiddleware(userConfig, ipConfig, globalConfig *routing.RateLimitConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		userLimiter:   routing.NewRateLimiter(userConfig),
		ipLimiter:     routing.NewRateLimiter(ipConfig),
		globalLimiter: routing.NewRateLimiter(globalConfig),
	}
}

func (m *RateLimitMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "default"
		}

		if !m.globalLimiter.Allow("global") {
			m.writeRateLimitError(w, "Global rate limit exceeded")
			return
		}

		if !m.ipLimiter.Allow(ip) {
			m.writeRateLimitError(w, "IP rate limit exceeded")
			return
		}

		if !m.userLimiter.Allow(userID) {
			m.writeRateLimitError(w, "User rate limit exceeded")
			return
		}

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", m.userLimiter.GetRemaining(userID)))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", m.userLimiter.GetRemaining(userID)))

		next.ServeHTTP(w, r)
	})
}

func (m *RateLimitMiddleware) writeRateLimitError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", "60")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(fmt.Sprintf(`{"error":{"message":"%s","type":"rate_limit_error"}}`, message)))
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}
