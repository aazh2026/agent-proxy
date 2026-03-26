package middleware

import (
	"net/http"
)

type SecurityHeadersMiddleware struct {
	hsts          bool
	csp           string
	frameOptions  string
	xssProtection string
}

func NewSecurityHeadersMiddleware() *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{
		hsts:          true,
		csp:           "default-src 'self'",
		frameOptions:  "DENY",
		xssProtection: "1; mode=block",
	}
}

func (m *SecurityHeadersMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.hsts {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		if m.csp != "" {
			w.Header().Set("Content-Security-Policy", m.csp)
		}
		if m.frameOptions != "" {
			w.Header().Set("X-Frame-Options", m.frameOptions)
		}
		if m.xssProtection != "" {
			w.Header().Set("X-XSS-Protection", m.xssProtection)
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
