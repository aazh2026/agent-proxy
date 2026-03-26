package auth

import (
	"net/http"
	"strings"
)

type SessionAuthenticator struct {
	sessionStore *SessionStore
}

func NewSessionAuthenticator(sessionStore *SessionStore) *SessionAuthenticator {
	return &SessionAuthenticator{
		sessionStore: sessionStore,
	}
}

func (a *SessionAuthenticator) Authenticate(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrInvalidSession
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidSession
	}

	token := parts[1]
	session, err := a.sessionStore.GetSession(token)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", ErrInvalidSession
	}

	return session.UserID, nil
}
