package testutil

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDB(t *testing.T) *sql.DB {
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
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			disabled BOOLEAN DEFAULT FALSE
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

func CreateTestUser(t *testing.T, db *sql.DB, userID, username, passwordHash string) {
	t.Helper()

	_, err := db.Exec(
		"INSERT INTO users (user_id, username, password_hash) VALUES (?, ?, ?)",
		userID, username, passwordHash,
	)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
}

func CreateTestToken(t *testing.T, db *sql.DB, tokenID, userID, provider string, expiresAt int64) {
	t.Helper()

	_, err := db.Exec(
		"INSERT INTO tokens (token_id, user_id, provider, type, expires_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))",
		tokenID, userID, provider, "bearer", expiresAt,
	)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}
}
