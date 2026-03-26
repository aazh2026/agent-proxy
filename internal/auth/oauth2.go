package auth

import (
	"net/http"
)

type OAuth2Authenticator struct {
	clientID     string
	clientSecret string
	redirectURL  string
	allowedOrgs  []string
}

func NewOAuth2Authenticator(clientID, clientSecret, redirectURL string, allowedOrgs []string) *OAuth2Authenticator {
	return &OAuth2Authenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		allowedOrgs:  allowedOrgs,
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
	http.Error(w, "OAuth2 not implemented", http.StatusNotImplemented)
}

func (h *OAuth2Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "OAuth2 not implemented", http.StatusNotImplemented)
}
