package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

type Session struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) CreateSession(userID string, ttlSeconds int) (*Session, error) {
	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(ttlSeconds) * time.Second)

	_, err = s.db.Exec(
		"INSERT INTO sessions (session_id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)",
		token, userID, expiresAt, now,
	)
	if err != nil {
		return nil, err
	}

	return &Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}, nil
}

func (s *SessionStore) GetSession(token string) (*Session, error) {
	session := &Session{}
	err := s.db.QueryRow(
		"SELECT session_id, user_id, expires_at, created_at FROM sessions WHERE session_id = ?",
		token,
	).Scan(&session.Token, &session.UserID, &session.ExpiresAt, &session.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(token)
		return nil, ErrSessionExpired
	}

	return session, nil
}

func (s *SessionStore) DeleteSession(token string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE session_id = ?", token)
	return err
}

func (s *SessionStore) CleanupExpiredSessions() error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
