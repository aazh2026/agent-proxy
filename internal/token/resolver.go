package token

import (
	"sync"
	"time"

	"github.com/openclaw/agent-proxy/internal/crypto"
)

type TokenResolver struct {
	tokenStore *TokenStore
	encryptor  *crypto.Encryptor
	cache      sync.Map
}

func NewTokenResolver(tokenStore *TokenStore, encryptor *crypto.Encryptor) *TokenResolver {
	return &TokenResolver{
		tokenStore: tokenStore,
		encryptor:  encryptor,
	}
}

func (r *TokenResolver) ResolveToken(userID, provider, model string) (*ResolvedToken, error) {
	tokens, err := r.tokenStore.GetTokensByProvider(userID, provider)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		if !r.isTokenAllowed(token, model) {
			continue
		}

		if token.ExpiresAt > 0 && isExpired(token.ExpiresAt) {
			continue
		}

		accessToken, err := r.decryptToken(token.AccessTokenEncrypted)
		if err != nil {
			continue
		}

		var refreshToken string
		if len(token.RefreshTokenEncrypted) > 0 {
			refreshToken, _ = r.decryptToken(token.RefreshTokenEncrypted)
		}

		return &ResolvedToken{
			TokenID:      token.TokenID,
			UserID:       token.UserID,
			Provider:     token.Provider,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    token.ExpiresAt,
		}, nil
	}

	return nil, ErrNoAvailableToken
}

func (r *TokenResolver) isTokenAllowed(token *Token, model string) bool {
	if len(token.AllowedModels) == 0 {
		return true
	}

	for _, allowed := range token.AllowedModels {
		if allowed == model {
			return true
		}
	}
	return false
}

func (r *TokenResolver) decryptToken(encrypted []byte) (string, error) {
	decrypted, err := r.encryptor.Decrypt(encrypted)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

type ResolvedToken struct {
	TokenID      string
	UserID       string
	Provider     string
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}

func (t *ResolvedToken) Clear() {
	t.AccessToken = ""
	t.RefreshToken = ""
}

func isExpired(expiresAt int64) bool {
	return time.Now().Unix() > expiresAt
}
