package token

import "errors"

var (
	ErrNoAvailableToken = errors.New("no available token for provider")
)
