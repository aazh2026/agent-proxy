package auth

import "errors"

var (
	ErrInvalidUserID   = errors.New("invalid user ID format")
	ErrUserNotAllowed  = errors.New("user not in allowed list")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserDisabled    = errors.New("user account is disabled")
	ErrSessionExpired  = errors.New("session expired")
	ErrInvalidSession  = errors.New("invalid session")
)
