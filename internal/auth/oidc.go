package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OIDCAuthenticator struct {
	config         *oauth2.Config
	allowedDomains []string
	stateStore     map[string]time.Time
}

func NewOIDCAuthenticator(clientID, clientSecret, redirectURL string, allowedDomains []string) *OIDCAuthenticator {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &OIDCAuthenticator{
		config:         config,
		allowedDomains: allowedDomains,
		stateStore:     make(map[string]time.Time),
	}
}

func (a *OIDCAuthenticator) Authenticate(r *http.Request) (string, error) {
	return "", ErrInvalidSession
}

type OIDCHandler struct {
	authenticator *OIDCAuthenticator
	sessionStore  *SessionStore
}

func NewOIDCHandler(authenticator *OIDCAuthenticator, sessionStore *SessionStore) *OIDCHandler {
	return &OIDCHandler{
		authenticator: authenticator,
		sessionStore:  sessionStore,
	}
}

func (h *OIDCHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	h.authenticator.stateStore[state] = time.Now().Add(10 * time.Minute)

	url := h.authenticator.config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *OIDCHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if state == "" || code == "" {
		http.Error(w, "Missing state or code", http.StatusBadRequest)
		return
	}

	expiry, ok := h.authenticator.stateStore[state]
	if !ok || time.Now().After(expiry) {
		http.Error(w, "Invalid or expired state", http.StatusBadRequest)
		return
	}
	delete(h.authenticator.stateStore, state)

	ctx := context.Background()
	token, err := h.authenticator.config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange code", http.StatusUnauthorized)
		return
	}

	client := h.authenticator.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		ID    string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	if len(h.authenticator.allowedDomains) > 0 {
		domain := extractDomain(userInfo.Email)
		if !isDomainAllowed(domain, h.authenticator.allowedDomains) {
			http.Error(w, "Domain not allowed", http.StatusForbidden)
			return
		}
	}

	session, err := h.sessionStore.CreateSession(userInfo.Email, 24*60*60)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      session.Token,
		"expires_in": 24 * 60 * 60,
		"user_id":    userInfo.Email,
	})
}

func generateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func isDomainAllowed(domain string, allowedDomains []string) bool {
	for _, allowed := range allowedDomains {
		if strings.EqualFold(domain, allowed) {
			return true
		}
	}
	return false
}

func (h *OIDCHandler) writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"type":    "authentication_error",
		},
	})
}
