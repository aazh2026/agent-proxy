package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OAuth2Authenticator struct {
	config      *oauth2.Config
	allowedOrgs []string
	stateStore  map[string]time.Time
}

func NewOAuth2Authenticator(clientID, clientSecret, redirectURL string, allowedOrgs []string) *OAuth2Authenticator {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}

	return &OAuth2Authenticator{
		config:      config,
		allowedOrgs: allowedOrgs,
		stateStore:  make(map[string]time.Time),
	}
}

func (a *OAuth2Authenticator) Authenticate(r *http.Request) (string, error) {
	return "", ErrInvalidSession
}

type OAuth2Handler struct {
	authenticator *OAuth2Authenticator
	sessionStore  *SessionStore
}

func NewOAuth2Handler(authenticator *OAuth2Authenticator, sessionStore *SessionStore) *OAuth2Handler {
	return &OAuth2Handler{
		authenticator: authenticator,
		sessionStore:  sessionStore,
	}
}

func (h *OAuth2Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
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

func (h *OAuth2Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
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

	userResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	var userInfo struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	if len(h.authenticator.allowedOrgs) > 0 {
		orgResp, err := client.Get("https://api.github.com/user/orgs")
		if err != nil {
			http.Error(w, "Failed to get organizations", http.StatusInternalServerError)
			return
		}
		defer orgResp.Body.Close()

		var orgs []struct {
			Login string `json:"login"`
		}
		if err := json.NewDecoder(orgResp.Body).Decode(&orgs); err != nil {
			http.Error(w, "Failed to decode organizations", http.StatusInternalServerError)
			return
		}

		if !isUserInOrgs(userInfo.Login, orgs, h.authenticator.allowedOrgs) {
			http.Error(w, "Organization not allowed", http.StatusForbidden)
			return
		}
	}

	userID := userInfo.Email
	if userID == "" {
		userID = userInfo.Login
	}

	session, err := h.sessionStore.CreateSession(userID, 24*60*60)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      session.Token,
		"expires_in": 24 * 60 * 60,
		"user_id":    userID,
	})
}

func isUserInOrgs(username string, orgs interface{}, allowedOrgs []string) bool {
	orgList, ok := orgs.([]struct{ Login string })
	if !ok {
		return false
	}

	orgMap := make(map[string]bool)
	for _, org := range orgList {
		orgMap[org.Login] = true
	}

	for _, allowed := range allowedOrgs {
		if orgMap[allowed] {
			return true
		}
	}
	return false
}
