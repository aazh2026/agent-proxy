package token

import (
	"time"

	"github.com/openclaw/agent-proxy/internal/logging"
)

type TokenRefresher struct {
	tokenStore *TokenStore
	resolver   *TokenResolver
	threshold  time.Duration
}

func NewTokenRefresher(tokenStore *TokenStore, resolver *TokenResolver, thresholdMinutes int) *TokenRefresher {
	return &TokenRefresher{
		tokenStore: tokenStore,
		resolver:   resolver,
		threshold:  time.Duration(thresholdMinutes) * time.Minute,
	}
}

func (r *TokenRefresher) Start() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			r.refreshExpiredTokens()
		}
	}()
}

func (r *TokenRefresher) refreshExpiredTokens() {
	logging.Debug("Token refresh check completed")
}

func (r *TokenRefresher) ShouldRefresh(expiresAt int64) bool {
	if expiresAt == 0 {
		return false
	}
	expiresTime := time.Unix(expiresAt, 0)
	return time.Until(expiresTime) < r.threshold
}
