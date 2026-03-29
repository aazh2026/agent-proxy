package routing

import (
	"errors"
	"testing"
	"time"

	"github.com/openclaw/agent-proxy/internal/token"
)

func TestFailoverHandler_ExecuteWithRetry_Success(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)
	handler := NewFailoverHandler(3, 100, 1000, selector)

	tokens := []*token.Token{
		{TokenID: "token1", Status: "enabled", Priority: 1},
		{TokenID: "token2", Status: "enabled", Priority: 1},
	}

	callCount := 0
	err := handler.ExecuteWithRetry(tokens, func(tok *token.Token) error {
		callCount++
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

func TestFailoverHandler_ExecuteWithRetry_AllFail(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)
	handler := NewFailoverHandler(3, 10, 100, selector)

	tokens := []*token.Token{
		{TokenID: "token1", Status: "enabled", Priority: 1},
		{TokenID: "token2", Status: "enabled", Priority: 1},
	}

	callCount := 0
	testErr := errors.New("test error")
	err := handler.ExecuteWithRetry(tokens, func(tok *token.Token) error {
		callCount++
		return testErr
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != testErr {
		t.Errorf("Expected test error, got %v", err)
	}

	if callCount != 4 {
		t.Errorf("Expected 4 calls (1 initial + 3 retries), got %d", callCount)
	}
}

func TestFailoverHandler_ExecuteWithRetry_Failover(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)
	handler := NewFailoverHandler(3, 10, 100, selector)

	tokens := []*token.Token{
		{TokenID: "token1", Status: "enabled", Priority: 1},
		{TokenID: "token2", Status: "enabled", Priority: 1},
	}

	attempts := 0
	err := handler.ExecuteWithRetry(tokens, func(tok *token.Token) error {
		attempts++
		if tok.TokenID == "token1" {
			return errors.New("token1 failed")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestFailoverHandler_ExecuteWithRetry_NoTokens(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)
	handler := NewFailoverHandler(3, 100, 1000, selector)

	err := handler.ExecuteWithRetry([]*token.Token{}, func(tok *token.Token) error {
		return nil
	})

	if err != ErrNoAvailableToken {
		t.Errorf("Expected ErrNoAvailableToken, got %v", err)
	}
}

func TestFailoverHandler_ExecuteWithRetry_DisabledTokens(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)
	handler := NewFailoverHandler(3, 10, 100, selector)

	tokens := []*token.Token{
		{TokenID: "token1", Status: "disabled", Priority: 1},
		{TokenID: "token2", Status: "disabled", Priority: 1},
	}

	err := handler.ExecuteWithRetry(tokens, func(tok *token.Token) error {
		return nil
	})

	if err == nil {
		t.Error("Expected error for disabled tokens, got nil")
	}
}

func TestFailoverHandler_CalculateDelay(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)
	handler := NewFailoverHandler(3, 100, 1000, selector)

	delay0 := handler.calculateDelay(0)
	if delay0 != 100*time.Millisecond {
		t.Errorf("Expected 100ms for attempt 0, got %v", delay0)
	}

	delay1 := handler.calculateDelay(1)
	if delay1 != 200*time.Millisecond {
		t.Errorf("Expected 200ms for attempt 1, got %v", delay1)
	}

	delay2 := handler.calculateDelay(2)
	if delay2 != 400*time.Millisecond {
		t.Errorf("Expected 400ms for attempt 2, got %v", delay2)
	}

	delay3 := handler.calculateDelay(3)
	if delay3 != 800*time.Millisecond {
		t.Errorf("Expected 800ms for attempt 3, got %v", delay3)
	}

	delay10 := handler.calculateDelay(10)
	if delay10 != 1000*time.Millisecond {
		t.Errorf("Expected 1000ms (max) for attempt 10, got %v", delay10)
	}
}
