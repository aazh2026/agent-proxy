package token

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE tokens (
			token_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			provider TEXT NOT NULL,
			type TEXT NOT NULL,
			access_token_encrypted BLOB NOT NULL,
			refresh_token_encrypted BLOB,
			expires_at INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT DEFAULT 'enabled',
			priority INTEGER DEFAULT 0,
			allowed_models TEXT DEFAULT '[]'
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return db
}

func TestTokenStore(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewTokenStore(db)

	token := &Token{
		TokenID:              "tk_test123",
		UserID:               "user1",
		Provider:             "openai",
		Type:                 "api_key",
		AccessTokenEncrypted: []byte("encrypted_data"),
		Status:               "enabled",
		Priority:             1,
		AllowedModels:        []string{"gpt-4"},
	}

	err := store.CreateToken(token)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	retrieved, err := store.GetToken("tk_test123")
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if retrieved == nil {
		t.Fatal("Token not found")
	}
	if retrieved.UserID != "user1" {
		t.Errorf("Expected user1, got %s", retrieved.UserID)
	}
}

func TestGetTokensByUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewTokenStore(db)

	tokens := []*Token{
		{
			TokenID:              "tk_1",
			UserID:               "user1",
			Provider:             "openai",
			Type:                 "api_key",
			AccessTokenEncrypted: []byte("data1"),
			Status:               "enabled",
		},
		{
			TokenID:              "tk_2",
			UserID:               "user1",
			Provider:             "anthropic",
			Type:                 "api_key",
			AccessTokenEncrypted: []byte("data2"),
			Status:               "enabled",
		},
		{
			TokenID:              "tk_3",
			UserID:               "user2",
			Provider:             "openai",
			Type:                 "api_key",
			AccessTokenEncrypted: []byte("data3"),
			Status:               "enabled",
		},
	}

	for _, token := range tokens {
		err := store.CreateToken(token)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}
	}

	user1Tokens, err := store.GetTokensByUser("user1")
	if err != nil {
		t.Fatalf("Failed to get tokens: %v", err)
	}
	if len(user1Tokens) != 2 {
		t.Errorf("Expected 2 tokens, got %d", len(user1Tokens))
	}
}

func TestUpdateTokenStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewTokenStore(db)

	token := &Token{
		TokenID:              "tk_test",
		UserID:               "user1",
		Provider:             "openai",
		Type:                 "api_key",
		AccessTokenEncrypted: []byte("data"),
		Status:               "enabled",
	}

	err := store.CreateToken(token)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	err = store.UpdateTokenStatus("tk_test", "disabled")
	if err != nil {
		t.Fatalf("Failed to update token status: %v", err)
	}

	retrieved, err := store.GetToken("tk_test")
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if retrieved.Status != "disabled" {
		t.Errorf("Expected disabled, got %s", retrieved.Status)
	}
}

func TestDeleteToken(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewTokenStore(db)

	token := &Token{
		TokenID:              "tk_test",
		UserID:               "user1",
		Provider:             "openai",
		Type:                 "api_key",
		AccessTokenEncrypted: []byte("data"),
		Status:               "enabled",
	}

	err := store.CreateToken(token)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	err = store.DeleteToken("tk_test")
	if err != nil {
		t.Fatalf("Failed to delete token: %v", err)
	}

	retrieved, err := store.GetToken("tk_test")
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if retrieved != nil {
		t.Error("Expected nil, got token")
	}
}
