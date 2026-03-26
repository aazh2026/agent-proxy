package auth

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuthMethod string

const (
	AuthMethodXUserID AuthMethod = "x-user-id"
	AuthMethodLocal   AuthMethod = "local"
	AuthMethodOIDC    AuthMethod = "oidc"
	AuthMethodOAuth2  AuthMethod = "oauth2"
	AuthMethodSession AuthMethod = "session"
)

type Authenticator interface {
	Authenticate(r *http.Request) (string, error)
}

type XUserIDAuthenticator struct {
	allowedUserIDs map[string]bool
	defaultUser    string
}

func NewXUserIDAuthenticator(allowedUserIDs []string) *XUserIDAuthenticator {
	allowed := make(map[string]bool)
	for _, id := range allowedUserIDs {
		allowed[id] = true
	}

	return &XUserIDAuthenticator{
		allowedUserIDs: allowed,
		defaultUser:    "default",
	}
}

func (a *XUserIDAuthenticator) Authenticate(r *http.Request) (string, error) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return a.defaultUser, nil
	}

	if !isValidUserID(userID) {
		return "", ErrInvalidUserID
	}

	if len(a.allowedUserIDs) > 0 && !a.allowedUserIDs[userID] {
		return "", ErrUserNotAllowed
	}

	return userID, nil
}

var userIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

func isValidUserID(userID string) bool {
	return userIDPattern.MatchString(userID)
}

func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

type AuthMiddleware struct {
	authenticator Authenticator
}

func NewAuthMiddleware(authenticator Authenticator) *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: authenticator,
	}
}

func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicEndpoint(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := m.authenticator.Authenticate(r)
		if err != nil {
			writeAuthError(w, err)
			return
		}

		ctx := SetUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/auth/login",
		"/auth/oidc/login",
		"/auth/oidc/callback",
		"/auth/oauth2/login",
		"/auth/oauth2/callback",
	}

	for _, p := range publicPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func writeAuthError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":{"message":"` + err.Error() + `","type":"authentication_error"}}`))
}
