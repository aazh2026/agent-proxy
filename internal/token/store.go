package token

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Token struct {
	TokenID               string    `json:"token_id"`
	UserID                string    `json:"user_id"`
	Provider              string    `json:"provider"`
	Type                  string    `json:"type"`
	AccessTokenEncrypted  []byte    `json:"-"`
	RefreshTokenEncrypted []byte    `json:"-"`
	ExpiresAt             int64     `json:"expires_at"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Status                string    `json:"status"`
	Priority              int       `json:"priority"`
	AllowedModels         []string  `json:"allowed_models"`
}

type TokenStore struct {
	db *sql.DB
}

func NewTokenStore(db *sql.DB) *TokenStore {
	return &TokenStore{db: db}
}

func (s *TokenStore) CreateToken(token *Token) error {
	allowedModelsJSON, err := json.Marshal(token.AllowedModels)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = s.db.Exec(
		`INSERT INTO tokens (
			token_id, user_id, provider, type,
			access_token_encrypted, refresh_token_encrypted,
			expires_at, created_at, updated_at,
			status, priority, allowed_models
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		token.TokenID, token.UserID, token.Provider, token.Type,
		token.AccessTokenEncrypted, token.RefreshTokenEncrypted,
		token.ExpiresAt, now, now,
		token.Status, token.Priority, string(allowedModelsJSON),
	)
	if err != nil {
		return err
	}

	token.CreatedAt = now
	token.UpdatedAt = now
	return nil
}

func (s *TokenStore) GetToken(tokenID string) (*Token, error) {
	token := &Token{}
	var allowedModelsJSON string
	err := s.db.QueryRow(
		`SELECT token_id, user_id, provider, type,
			access_token_encrypted, refresh_token_encrypted,
			expires_at, created_at, updated_at,
			status, priority, allowed_models
		FROM tokens WHERE token_id = ?`,
		tokenID,
	).Scan(
		&token.TokenID, &token.UserID, &token.Provider, &token.Type,
		&token.AccessTokenEncrypted, &token.RefreshTokenEncrypted,
		&token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt,
		&token.Status, &token.Priority, &allowedModelsJSON,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(allowedModelsJSON), &token.AllowedModels)
	return token, nil
}

func (s *TokenStore) GetTokensByUser(userID string) ([]*Token, error) {
	rows, err := s.db.Query(
		`SELECT token_id, user_id, provider, type,
			expires_at, created_at, updated_at,
			status, priority, allowed_models
		FROM tokens WHERE user_id = ? ORDER BY priority DESC, created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*Token
	for rows.Next() {
		token := &Token{}
		var allowedModelsJSON string
		err := rows.Scan(
			&token.TokenID, &token.UserID, &token.Provider, &token.Type,
			&token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt,
			&token.Status, &token.Priority, &allowedModelsJSON,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(allowedModelsJSON), &token.AllowedModels)
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (s *TokenStore) GetAllTokens() ([]*Token, error) {
	rows, err := s.db.Query(
		`SELECT token_id, user_id, provider, type,
			expires_at, created_at, updated_at,
			status, priority, allowed_models
		FROM tokens WHERE status = 'active' ORDER BY priority DESC, created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*Token
	for rows.Next() {
		token := &Token{}
		var allowedModelsJSON string
		err := rows.Scan(
			&token.TokenID, &token.UserID, &token.Provider, &token.Type,
			&token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt,
			&token.Status, &token.Priority, &allowedModelsJSON,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(allowedModelsJSON), &token.AllowedModels)
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (s *TokenStore) GetTokensByProvider(userID, provider string) ([]*Token, error) {
	rows, err := s.db.Query(
		`SELECT token_id, user_id, provider, type,
			expires_at, created_at, updated_at,
			status, priority, allowed_models
		FROM tokens WHERE user_id = ? AND provider = ? AND status = 'enabled'
		ORDER BY priority DESC`,
		userID, provider,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*Token
	for rows.Next() {
		token := &Token{}
		var allowedModelsJSON string
		err := rows.Scan(
			&token.TokenID, &token.UserID, &token.Provider, &token.Type,
			&token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt,
			&token.Status, &token.Priority, &allowedModelsJSON,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(allowedModelsJSON), &token.AllowedModels)
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (s *TokenStore) UpdateToken(token *Token) error {
	allowedModelsJSON, err := json.Marshal(token.AllowedModels)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		`UPDATE tokens SET
			provider = ?, type = ?,
			access_token_encrypted = ?, refresh_token_encrypted = ?,
			expires_at = ?, updated_at = ?,
			status = ?, priority = ?, allowed_models = ?
		WHERE token_id = ?`,
		token.Provider, token.Type,
		token.AccessTokenEncrypted, token.RefreshTokenEncrypted,
		token.ExpiresAt, time.Now(),
		token.Status, token.Priority, string(allowedModelsJSON),
		token.TokenID,
	)
	return err
}

func (s *TokenStore) UpdateTokenStatus(tokenID, status string) error {
	_, err := s.db.Exec(
		"UPDATE tokens SET status = ?, updated_at = ? WHERE token_id = ?",
		status, time.Now(), tokenID,
	)
	return err
}

func (s *TokenStore) DeleteToken(tokenID string) error {
	_, err := s.db.Exec("DELETE FROM tokens WHERE token_id = ?", tokenID)
	return err
}

func (s *TokenStore) DeleteTokensByUser(userID string) error {
	_, err := s.db.Exec("DELETE FROM tokens WHERE user_id = ?", userID)
	return err
}
