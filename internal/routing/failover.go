package routing

import (
	"time"

	"github.com/openclaw/agent-proxy/internal/token"
)

type FailoverHandler struct {
	maxRetries int
	retryDelay time.Duration
	maxDelay   time.Duration
	selector   *TokenSelector
}

func NewFailoverHandler(maxRetries int, retryDelayMs, maxDelayMs int, selector *TokenSelector) *FailoverHandler {
	return &FailoverHandler{
		maxRetries: maxRetries,
		retryDelay: time.Duration(retryDelayMs) * time.Millisecond,
		maxDelay:   time.Duration(maxDelayMs) * time.Millisecond,
		selector:   selector,
	}
}

func (h *FailoverHandler) ExecuteWithRetry(tokens []*token.Token, fn func(*token.Token) error) error {
	if len(tokens) == 0 {
		return ErrNoAvailableToken
	}

	var lastErr error
	usedTokens := make(map[string]bool)

	for attempt := 0; attempt <= h.maxRetries; attempt++ {
		availableTokens := h.filterAvailable(tokens, usedTokens)
		if len(availableTokens) == 0 {
			// If we've tried all tokens, reset usedTokens to allow cycling
			if len(usedTokens) == len(tokens) {
				usedTokens = make(map[string]bool)
				availableTokens = h.filterAvailable(tokens, usedTokens)
				if len(availableTokens) == 0 {
					return ErrNoAvailableToken
				}
			} else {
				return ErrNoAvailableToken
			}
		}

		selectedToken := h.selector.SelectToken(availableTokens)
		if selectedToken == nil {
			return ErrNoAvailableToken
		}

		err := fn(selectedToken)
		if err == nil {
			return nil
		}

		lastErr = err
		usedTokens[selectedToken.TokenID] = true

		if attempt < h.maxRetries {
			delay := h.calculateDelay(attempt)
			time.Sleep(delay)
		}
	}

	if lastErr == nil {
		return ErrNoAvailableToken
	}
	return lastErr
}

func (h *FailoverHandler) filterAvailable(tokens []*token.Token, used map[string]bool) []*token.Token {
	var available []*token.Token
	for _, t := range tokens {
		if !used[t.TokenID] && t.Status == "enabled" {
			available = append(available, t)
		}
	}
	return available
}

func (h *FailoverHandler) calculateDelay(attempt int) time.Duration {
	delay := h.retryDelay * time.Duration(1<<uint(attempt))
	if delay > h.maxDelay {
		delay = h.maxDelay
	}
	return delay
}
