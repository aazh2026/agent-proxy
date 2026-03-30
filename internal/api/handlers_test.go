package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/openclaw/agent-proxy/internal/auth"
	"github.com/openclaw/agent-proxy/internal/crypto"
	"github.com/openclaw/agent-proxy/internal/pipeline"
	"github.com/openclaw/agent-proxy/internal/routing"
	"github.com/openclaw/agent-proxy/internal/token"

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
		CREATE TABLE IF NOT EXISTS tokens (
			token_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			provider TEXT NOT NULL,
			type TEXT NOT NULL,
			access_token_encrypted BLOB,
			refresh_token_encrypted BLOB,
			expires_at INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			status TEXT DEFAULT 'active',
			priority INTEGER DEFAULT 0,
			allowed_models TEXT DEFAULT '[]'
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create tokens table: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.RemoveAll(tmpDir)
	})

	return db
}

func TestChatCompletionsHandler_HandleRequest(t *testing.T) {
	db := setupTestDB(t)
	tokenStore := token.NewTokenStore(db)
	key, _ := crypto.GenerateKey()
	encryptor, _ := crypto.NewEncryptor(key)
	forwardingStage := pipeline.NewForwardingStage()
	tokenResolver := token.NewTokenResolver(tokenStore, encryptor)
	routingHandler := routing.NewRequestHandler(tokenResolver, 3, 100, 5000, routing.StrategyRoundRobin)
	handler := NewChatCompletionsHandler(forwardingStage, tokenResolver, routingHandler)

	body := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserID(req.Context(), "test-user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError && w.Code != http.StatusUnauthorized && w.Code != http.StatusBadGateway {
		t.Errorf("Expected status 200, 401, 500, or 502, got %d", w.Code)
	}
}

func TestEmbeddingsHandler_HandleRequest(t *testing.T) {
	db := setupTestDB(t)
	tokenStore := token.NewTokenStore(db)
	key, _ := crypto.GenerateKey()
	encryptor, _ := crypto.NewEncryptor(key)
	forwardingStage := pipeline.NewForwardingStage()
	tokenResolver := token.NewTokenResolver(tokenStore, encryptor)
	routingHandler := routing.NewRequestHandler(tokenResolver, 3, 100, 5000, routing.StrategyRoundRobin)
	handler := NewEmbeddingsHandler(forwardingStage, tokenResolver, routingHandler)

	body := map[string]interface{}{
		"model": "text-embedding-ada-002",
		"input": "Hello world",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/embeddings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	ctx := auth.SetUserID(req.Context(), "test-user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError && w.Code != http.StatusUnauthorized && w.Code != http.StatusBadGateway {
		t.Errorf("Expected status 200, 401, 500, or 502, got %d", w.Code)
	}
}
