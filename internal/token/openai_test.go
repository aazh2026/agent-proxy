package token

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/openclaw/agent-proxy/internal/auth"
	"github.com/openclaw/agent-proxy/internal/crypto"
	"github.com/openclaw/agent-proxy/internal/db"
)

func setupOpenAIOAuthTestDB(t *testing.T) (*db.DB, func()) {
	tmpDir, err := os.MkdirTemp("", "agent-proxy-test-*")
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	dbPath := filepath.Join(tmpDir, "agent-proxy.db")
	database, err := db.New(dbPath)
	if err != nil {
		cleanup()
		t.Fatalf("failed to open db: %v", err)
	}

	return database, cleanup
}

func TestOpenAIOAuthHandler(t *testing.T) {
	database, cleanup := setupOpenAIOAuthTestDB(t)
	defer cleanup()

	userStore := auth.NewUserStore(database.Conn())
	_, err := userStore.CreateUser("alice", "alice", "password", "alice@example.com")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	encryptor, err := crypto.NewEncryptor(make([]byte, 32))
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	// Use temporary endpoints for test.
	mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		if r.FormValue("grant_type") != "authorization_code" {
			http.Error(w, "bad grant_type", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  "openai-test-token",
			"refresh_token": "openai-test-refresh",
			"expires_in":    3600,
		})
	}))
	defer mockTokenServer.Close()

	tokenStore := NewTokenStore(database.Conn())
	h := NewOpenAIOAuthHandler("clientid", "clientsecret", "http://localhost/auth/openai/callback", tokenStore, encryptor)
	h.tokenEndpoint = mockTokenServer.URL
	h.authURL = "http://localhost/oauth/authorize"

	loginReq := httptest.NewRequest(http.MethodGet, "/auth/openai/login", nil)
	loginReq.Header.Set("X-User-ID", "alice")
	loginRes := httptest.NewRecorder()
	h.HandleLogin(loginRes, loginReq)

	if loginRes.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected redirect, got %d", loginRes.Code)
	}

	redirectURL := loginRes.Header().Get("Location")
	if redirectURL == "" {
		t.Fatal("expected redirect URL")
	}

	parsed, err := url.Parse(redirectURL)
	if err != nil {
		t.Fatal(err)
	}

	state := parsed.Query().Get("state")
	if state == "" {
		t.Fatal("state missing")
	}

	callbackReq := httptest.NewRequest(http.MethodGet, "/auth/openai/callback?code=testcode&state="+url.QueryEscape(state), nil)
	callbackRes := httptest.NewRecorder()
	h.HandleCallback(callbackRes, callbackReq)

	if callbackRes.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", callbackRes.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(callbackRes.Body).Decode(&result); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if result["status"] != "ok" {
		t.Fatalf("expected ok status, got %v", result["status"])
	}

	tokens, err := tokenStore.GetTokensByUser("alice")
	if err != nil {
		t.Fatalf("failed to list tokens: %v", err)
	}
	if len(tokens) == 0 {
		t.Fatal("expected at least one token")
	}

	retrieved, err := encryptor.Decrypt(tokens[0].AccessTokenEncrypted)
	if err != nil {
		t.Fatalf("failed to decrypt token: %v", err)
	}
	if string(retrieved) != "openai-test-token" {
		t.Fatalf("unexpected access token %s", string(retrieved))
	}
}
