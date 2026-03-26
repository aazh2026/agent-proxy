package routing

import "errors"

var (
	ErrNoAvailableToken = errors.New("no available token")
	ErrAllTokensFailed  = errors.New("all tokens failed")
	ErrQuotaExceeded    = errors.New("quota exceeded")
)
