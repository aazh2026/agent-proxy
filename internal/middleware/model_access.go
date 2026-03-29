package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/openclaw/agent-proxy/internal/auth"
)

type ModelAccessMiddleware struct {
	userModels    map[string][]string
	defaultModels []string
}

func NewModelAccessMiddleware(userModels map[string][]string, defaultModels []string) *ModelAccessMiddleware {
	return &ModelAccessMiddleware{
		userModels:    userModels,
		defaultModels: defaultModels,
	}
}

func (m *ModelAccessMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" && r.URL.Path != "/v1/embeddings" {
			next.ServeHTTP(w, r)
			return
		}

		userID := auth.GetUserID(r.Context())
		if userID == "" {
			userID = "default"
		}

		model := extractModelFromRequest(r)
		if model == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !m.isModelAllowed(userID, model) {
			m.writeAccessDenied(w, model)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *ModelAccessMiddleware) isModelAllowed(userID, model string) bool {
	allowedModels, ok := m.userModels[userID]
	if !ok {
		allowedModels = m.defaultModels
	}

	if len(allowedModels) == 0 {
		return true
	}

	for _, allowed := range allowedModels {
		if strings.HasPrefix(allowed, "*") {
			suffix := allowed[1:]
			if strings.HasSuffix(model, suffix) {
				return true
			}
		} else if allowed == model {
			return true
		}
	}

	return false
}

func extractModelFromRequest(r *http.Request) string {
	if r.Method != http.MethodPost {
		return ""
	}

	var body struct {
		Model string `json:"model"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return ""
	}

	return body.Model
}

func (m *ModelAccessMiddleware) writeAccessDenied(w http.ResponseWriter, model string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": "Access denied for model: " + model,
			"type":    "access_denied",
		},
	})
}
