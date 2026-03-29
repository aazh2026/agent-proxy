package token

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/openclaw/agent-proxy/internal/crypto"

	_ "github.com/mattn/go-sqlite3"
)

func setupTokenTestDB(t *testing.T) *sql.DB {
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

func TestTokenRefresher_ShouldRefresh(t *testing.T) {
	db := setupTokenTestDB(t)
	store := NewTokenStore(db)
	key, _ := crypto.GenerateKey()
	encryptor, _ := crypto.NewEncryptor(key)
	resolver := NewTokenResolver(store, encryptor)
	refresher := NewTokenRefresher(store, resolver, 5)

	tests := []struct {
		name      string
		expiresAt int64
		expected  bool
	}{
		{
			name:      "Token expires in 1 minute (should refresh)",
			expiresAt: time.Now().Add(1 * time.Minute).Unix(),
			expected:  true,
		},
		{
			name:      "Token expires in 10 minutes (should not refresh)",
			expiresAt: time.Now().Add(10 * time.Minute).Unix(),
			expected:  false,
		},
		{
			name:      "Token expires at 0 (should not refresh)",
			expiresAt: 0,
			expected:  false,
		},
		{
			name:      "Token already expired (should refresh)",
			expiresAt: time.Now().Add(-1 * time.Minute).Unix(),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := refresher.ShouldRefresh(tt.expiresAt)
			if result != tt.expected {
				t.Errorf("ShouldRefresh(%d) = %v, expected %v", tt.expiresAt, result, tt.expected)
			}
		})
	}
}

func TestTokenRefresher_RefreshToken_NoRefreshToken(t *testing.T) {
	db := setupTokenTestDB(t)
	store := NewTokenStore(db)
	key, _ := crypto.GenerateKey()
	encryptor, _ := crypto.NewEncryptor(key)
	resolver := NewTokenResolver(store, encryptor)
	refresher := NewTokenRefresher(store, resolver, 5)

	token := &Token{
		TokenID:  "test-token",
		Provider: "openai",
		Status:   "active",
	}

	err := refresher.RefreshToken(token)
	if err == nil {
		t.Error("Expected error when no refresh token available, got nil")
	}
}

func TestTokenRefresher_RefreshToken_UnsupportedProvider(t *testing.T) {
	db := setupTokenTestDB(t)
	store := NewTokenStore(db)
	key, _ := crypto.GenerateKey()
	encryptor, _ := crypto.NewEncryptor(key)
	resolver := NewTokenResolver(store, encryptor)
	refresher := NewTokenRefresher(store, resolver, 5)

	token := &Token{
		TokenID:               "test-token",
		Provider:              "unsupported",
		RefreshTokenEncrypted: []byte("test-refresh-token"),
		Status:                "active",
	}

	err := refresher.RefreshToken(token)
	if err == nil {
		t.Error("Expected error for unsupported provider, got nil")
	}
}

func TestTokenRefresher_DisableToken(t *testing.T) {
	db := setupTokenTestDB(t)
	store := NewTokenStore(db)
	key, _ := crypto.GenerateKey()
	encryptor, _ := crypto.NewEncryptor(key)
	resolver := NewTokenResolver(store, encryptor)
	refresher := NewTokenRefresher(store, resolver, 5)

	token := &Token{
		TokenID:  "test-token",
		Provider: "openai",
		Status:   "active",
	}

	err := refresher.disableToken(token)
	if err != nil {
		t.Errorf("disableToken failed: %v", err)
	}

	if token.Status != "disabled" {
		t.Errorf("Expected token status to be 'disabled', got '%s'", token.Status)
	}
}

func TestTokenLifecycle_CreateAndRefresh(t *testing.T) {
	db := setupTokenTestDB(t)
	store := NewTokenStore(db)
	key, _ := crypto.GenerateKey()
	encryptor, _ := crypto.NewEncryptor(key)
	resolver := NewTokenResolver(store, encryptor)
	refresher := NewTokenRefresher(store, resolver, 5)

	token := &Token{
		TokenID:   "test-token",
		UserID:    "test-user",
		Provider:  "openai",
		Type:      "bearer",
		Status:    "active",
		Priority:  1,
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	}

	err := store.CreateToken(token)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	retrievedToken, err := store.GetToken("test-token")
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if retrievedToken == nil {
		t.Fatal("Token not found after creation")
	}

	if refresher.ShouldRefresh(retrievedToken.ExpiresAt) {
		t.Error("Token should not need refresh yet")
	}

	retrievedToken.ExpiresAt = time.Now().Add(1 * time.Minute).Unix()
	if !refresher.ShouldRefresh(retrievedToken.ExpiresAt) {
		t.Error("Token should need refresh")
	}
}
