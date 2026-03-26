package token

import (
	"database/sql"
	"time"
)

type AuditLog struct {
	ID        int64     `json:"id"`
	TokenID   string    `json:"token_id"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditStore struct {
	db *sql.DB
}

func NewAuditStore(db *sql.DB) *AuditStore {
	return &AuditStore{db: db}
}

func (s *AuditStore) Log(action, tokenID, userID, details, ip, userAgent string) error {
	_, err := s.db.Exec(
		`INSERT INTO audit_logs (token_id, user_id, action, details, ip, user_agent, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		tokenID, userID, action, details, ip, userAgent, time.Now(),
	)
	return err
}

func (s *AuditStore) GetLogsByToken(tokenID string, limit int) ([]*AuditLog, error) {
	rows, err := s.db.Query(
		`SELECT id, token_id, user_id, action, details, ip, user_agent, created_at
		FROM audit_logs WHERE token_id = ? ORDER BY created_at DESC LIMIT ?`,
		tokenID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*AuditLog
	for rows.Next() {
		log := &AuditLog{}
		err := rows.Scan(&log.ID, &log.TokenID, &log.UserID, &log.Action,
			&log.Details, &log.IP, &log.UserAgent, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (s *AuditStore) GetLogsByUser(userID string, limit int) ([]*AuditLog, error) {
	rows, err := s.db.Query(
		`SELECT id, token_id, user_id, action, details, ip, user_agent, created_at
		FROM audit_logs WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*AuditLog
	for rows.Next() {
		log := &AuditLog{}
		err := rows.Scan(&log.ID, &log.TokenID, &log.UserID, &log.Action,
			&log.Details, &log.IP, &log.UserAgent, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
