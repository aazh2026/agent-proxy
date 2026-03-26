package auth

import (
	"net/http"
)

type OIDCAuthenticator struct {
	clientID       string
	clientSecret   string
	redirectURL    string
	allowedDomains []string
}

func NewOIDCAuthenticator(clientID, clientSecret, redirectURL string, allowedDomains []string) *OIDCAuthenticator {
	return &OIDCAuthenticator{
		clientID:       clientID,
		clientSecret:   clientSecret,
		redirectURL:    redirectURL,
		allowedDomains: allowedDomains,
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
	http.Error(w, "OIDC not implemented", http.StatusNotImplemented)
}

func (h *OIDCHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "OIDC not implemented", http.StatusNotImplemented)
}
