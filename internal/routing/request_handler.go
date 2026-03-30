package routing

import (
	"context"
	"fmt"

	"github.com/openclaw/agent-proxy/internal/logging"
	"github.com/openclaw/agent-proxy/internal/token"
)

type RequestHandler struct {
	tokenResolver   *token.TokenResolver
	failoverHandler *FailoverHandler
	fallbackRouter  *FallbackRouter
	selector        *TokenSelector
}

func NewRequestHandler(
	tokenResolver *token.TokenResolver,
	maxRetries int,
	retryDelayMs int,
	maxDelayMs int,
	tokenStrategy TokenStrategy,
) *RequestHandler {
	selector := NewTokenSelector(tokenStrategy)
	failoverHandler := NewFailoverHandler(maxRetries, retryDelayMs, maxDelayMs, selector)
	fallbackRouter := NewFallbackRouter()

	return &RequestHandler{
		tokenResolver:   tokenResolver,
		failoverHandler: failoverHandler,
		fallbackRouter:  fallbackRouter,
		selector:        selector,
	}
}

func (h *RequestHandler) AddFallbackChain(model string, primary string, fallbacks []string) {
	h.fallbackRouter.AddChain(model, primary, fallbacks)
}

type TokenCallback func(provider string, token *token.Token) error

func (h *RequestHandler) ExecuteWithFailover(ctx context.Context, userID, model, primaryProvider string, callback TokenCallback) error {
	providers := []string{primaryProvider}

	if chain := h.fallbackRouter.GetChain(model); chain != nil {
		providers = chain.GetChain()
	}

	var lastErr error
	for _, provider := range providers {
		err := h.executeWithProviderRetry(ctx, userID, provider, model, callback)
		if err == nil {
			return nil
		}

		lastErr = err
		logging.Warn("Provider %s failed for model %s: %v", provider, model, err)
	}

	if lastErr == nil {
		return ErrAllTokensFailed
	}
	return lastErr
}

func (h *RequestHandler) executeWithProviderRetry(ctx context.Context, userID, provider, model string, callback TokenCallback) error {
	tokens, err := h.tokenResolver.GetValidTokens(userID, provider)
	if err != nil {
		return fmt.Errorf("failed to get tokens for provider %s: %w", provider, err)
	}

	if len(tokens) == 0 {
		return ErrNoAvailableToken
	}

	filteredTokens := h.filterByModel(tokens, model)
	if len(filteredTokens) == 0 {
		return ErrNoAvailableToken
	}

	return h.failoverHandler.ExecuteWithRetry(filteredTokens, func(t *token.Token) error {
		return callback(provider, t)
	})
}

func (h *RequestHandler) filterByModel(tokens []*token.Token, model string) []*token.Token {
	if model == "" {
		return tokens
	}

	var filtered []*token.Token
	for _, t := range tokens {
		if len(t.AllowedModels) == 0 {
			filtered = append(filtered, t)
			continue
		}
		for _, allowed := range t.AllowedModels {
			if allowed == model {
				filtered = append(filtered, t)
				break
			}
		}
	}
	return filtered
}
