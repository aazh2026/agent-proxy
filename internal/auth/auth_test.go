package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXUserIDAuthenticator(t *testing.T) {
	authenticator := NewXUserIDAuthenticator([]string{"alice", "bob"})

	tests := []struct {
		name      string
		userID    string
		expectErr bool
	}{
		{"Valid user", "alice", false},
		{"Valid user 2", "bob", false},
		{"Invalid user", "charlie", true},
		{"Empty user", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if test.userID != "" {
				req.Header.Set("X-User-ID", test.userID)
			}

			userID, err := authenticator.Authenticate(req)
			if test.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if test.userID == "" && userID != "default" {
					t.Errorf("Expected default user, got %s", userID)
				} else if test.userID != "" && userID != test.userID {
					t.Errorf("Expected %s, got %s", test.userID, userID)
				}
			}
		})
	}
}

func TestXUserIDAuthenticatorNoWhitelist(t *testing.T) {
	authenticator := NewXUserIDAuthenticator([]string{})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-User-ID", "anyuser")

	userID, err := authenticator.Authenticate(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if userID != "anyuser" {
		t.Errorf("Expected anyuser, got %s", userID)
	}
}

func TestAuthMiddleware(t *testing.T) {
	authenticator := NewXUserIDAuthenticator([]string{"alice"})
	middleware := NewAuthMiddleware(authenticator)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		w.Write([]byte(userID))
	}))

	req := httptest.NewRequest("GET", "/v1/chat/completions", nil)
	req.Header.Set("X-User-ID", "alice")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Body.String() != "alice" {
		t.Errorf("Expected alice, got %s", w.Body.String())
	}
}

func TestAuthMiddlewarePublicEndpoint(t *testing.T) {
	authenticator := NewXUserIDAuthenticator([]string{"alice"})
	middleware := NewAuthMiddleware(authenticator)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetUserID(t *testing.T) {
	ctx := SetUserID(nil, "test-user")
	userID := GetUserID(ctx)
	if userID != "test-user" {
		t.Errorf("Expected test-user, got %s", userID)
	}
}
