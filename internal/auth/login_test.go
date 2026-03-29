package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			session_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			expires_at DATETIME,
			created_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create sessions table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			email TEXT,
			enabled BOOLEAN DEFAULT TRUE,
			created_at DATETIME,
			updated_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.RemoveAll(tmpDir)
	})

	return db
}

func TestLoginHandler_LoginSuccess(t *testing.T) {
	db := setupTestDB(t)
	userStore := NewUserStore(db)
	sessionStore := NewSessionStore(db)
	handler := NewLoginHandler(userStore, sessionStore)

	_, err := userStore.CreateUser("test-user", "admin", "password", "admin@example.com")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	body := strings.NewReader(`{"username":"admin","password":"password"}`)
	req := httptest.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestLoginHandler_LoginFailure(t *testing.T) {
	db := setupTestDB(t)
	userStore := NewUserStore(db)
	sessionStore := NewSessionStore(db)
	handler := NewLoginHandler(userStore, sessionStore)

	body := strings.NewReader(`{"username":"admin","password":"wrong"}`)
	req := httptest.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestSessionStore_CreateAndGetSession(t *testing.T) {
	db := setupTestDB(t)
	sessionStore := NewSessionStore(db)

	session, err := sessionStore.CreateSession("test-user", 3600)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.Token == "" {
		t.Error("Session token should not be empty")
	}

	if session.UserID != "test-user" {
		t.Errorf("Expected user ID 'test-user', got '%s'", session.UserID)
	}

	retrievedSession, err := sessionStore.GetSession(session.Token)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrievedSession == nil {
		t.Fatal("Session not found")
	}

	if retrievedSession.UserID != "test-user" {
		t.Errorf("Expected user ID 'test-user', got '%s'", retrievedSession.UserID)
	}
}

func TestXUserIDAuthenticator_AccessControl(t *testing.T) {
	authenticator := NewXUserIDAuthenticator([]string{"admin", "user1"})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-User-ID", "admin")
	userID, err := authenticator.Authenticate(req)
	if err != nil {
		t.Errorf("Expected no error for allowed user, got %v", err)
	}
	if userID != "admin" {
		t.Errorf("Expected user ID 'admin', got '%s'", userID)
	}

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-User-ID", "unauthorized")
	_, err = authenticator.Authenticate(req)
	if err == nil {
		t.Error("Expected error for unauthorized user, got nil")
	}

	req = httptest.NewRequest("GET", "/", nil)
	userID, err = authenticator.Authenticate(req)
	if err != nil {
		t.Errorf("Expected no error for default user, got %v", err)
	}
	if userID != "default" {
		t.Errorf("Expected default user ID, got '%s'", userID)
	}
}

func TestAuthMiddleware_PublicEndpoints(t *testing.T) {
	authenticator := NewXUserIDAuthenticator([]string{})
	middleware := NewAuthMiddleware(authenticator)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	publicEndpoints := []string{"/health", "/auth/login", "/auth/oidc/login"}
	for _, endpoint := range publicEndpoints {
		req := httptest.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for public endpoint %s, got %d", endpoint, w.Code)
		}
	}
}
