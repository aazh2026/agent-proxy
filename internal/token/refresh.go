package token

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/openclaw/agent-proxy/internal/logging"
)

type TokenRefresher struct {
	tokenStore *TokenStore
	resolver   *TokenResolver
	threshold  time.Duration
	maxRetries int
	httpClient *http.Client
}

func NewTokenRefresher(tokenStore *TokenStore, resolver *TokenResolver, thresholdMinutes int) *TokenRefresher {
	return &TokenRefresher{
		tokenStore: tokenStore,
		resolver:   resolver,
		threshold:  time.Duration(thresholdMinutes) * time.Minute,
		maxRetries: 3,
		httpClient: &http.Client{Timeout: 30 * time.Second},
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
	logging.Debug("Checking for tokens that need refresh")

	tokens, err := r.tokenStore.GetAllTokens()
	if err != nil {
		logging.Error("Failed to get tokens for refresh check: %v", err)
		return
	}

	for _, token := range tokens {
		if token.Status != "active" {
			continue
		}

		if r.ShouldRefresh(token.ExpiresAt) {
			logging.Info("Token %s needs refresh (expires at %d)", token.TokenID, token.ExpiresAt)
			if err := r.RefreshToken(token); err != nil {
				logging.Error("Failed to refresh token %s: %v", token.TokenID, err)
			}
		}
	}
}

func (r *TokenRefresher) ShouldRefresh(expiresAt int64) bool {
	if expiresAt == 0 {
		return false
	}
	expiresTime := time.Unix(expiresAt, 0)
	return time.Until(expiresTime) < r.threshold
}

func (r *TokenRefresher) RefreshToken(token *Token) error {
	if token.RefreshTokenEncrypted == nil {
		return fmt.Errorf("no refresh token available")
	}

	for attempt := 0; attempt < r.maxRetries; attempt++ {
		err := r.attemptRefresh(token)
		if err == nil {
			logging.Info("Token refresh successful for token %s", token.TokenID)
			return nil
		}

		logging.Warn("Token refresh attempt %d failed for token %s: %v", attempt+1, token.TokenID, err)
		if attempt < r.maxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	logging.Error("Token refresh failed after %d attempts for token %s", r.maxRetries, token.TokenID)
	r.disableToken(token)
	return fmt.Errorf("token refresh failed after %d attempts", r.maxRetries)
}

func (r *TokenRefresher) attemptRefresh(token *Token) error {
	refreshToken, err := r.resolver.decryptToken(token.RefreshTokenEncrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt refresh token: %w", err)
	}
	defer func() { refreshToken = "" }()

	var refreshURL string
	switch token.Provider {
	case "openai":
		return fmt.Errorf("OpenAI API keys do not support refresh")
	case "anthropic":
		return fmt.Errorf("Anthropic API keys do not support refresh")
	case "google":
		refreshURL = "https://oauth2.googleapis.com/token"
	default:
		return fmt.Errorf("unsupported provider for refresh: %s", token.Provider)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", refreshURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Set("grant_type", "refresh_token")
	q.Set("refresh_token", refreshToken)
	q.Set("client_id", "")
	q.Set("client_secret", "")
	req.URL.RawQuery = q.Encode()

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send refresh request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh request failed with status %d", resp.StatusCode)
	}

	var refreshResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return fmt.Errorf("failed to decode refresh response: %w", err)
	}

	logging.Info("Token refresh successful, new token expires in %d seconds", refreshResp.ExpiresIn)

	return nil
}

func (r *TokenRefresher) disableToken(token *Token) error {
	token.Status = "disabled"
	if err := r.tokenStore.UpdateTokenStatus(token.TokenID, "disabled"); err != nil {
		return fmt.Errorf("failed to disable token: %w", err)
	}
	logging.Info("Token %s disabled after refresh failure", token.TokenID)
	return nil
}
