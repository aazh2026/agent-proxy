package routing

import (
	"testing"

	"github.com/openclaw/agent-proxy/internal/token"
)

func TestTokenSelector_RoundRobin(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)

	tokens := []*token.Token{
		{TokenID: "token1", Priority: 1},
		{TokenID: "token2", Priority: 1},
		{TokenID: "token3", Priority: 1},
	}

	selected := selector.SelectToken(tokens)
	if selected == nil {
		t.Fatal("Expected token, got nil")
	}

	selected2 := selector.SelectToken(tokens)
	if selected2 == nil {
		t.Fatal("Expected token, got nil")
	}

	selected3 := selector.SelectToken(tokens)
	if selected3 == nil {
		t.Fatal("Expected token, got nil")
	}

	selected4 := selector.SelectToken(tokens)
	if selected4 == nil {
		t.Fatal("Expected token, got nil")
	}
}

func TestTokenSelector_Weighted(t *testing.T) {
	selector := NewTokenSelector(StrategyWeighted)

	tokens := []*token.Token{
		{TokenID: "token1", Priority: 1},
		{TokenID: "token2", Priority: 2},
		{TokenID: "token3", Priority: 3},
	}

	selected := selector.SelectToken(tokens)
	if selected == nil {
		t.Fatal("Expected token, got nil")
	}
}

func TestTokenSelector_Priority(t *testing.T) {
	selector := NewTokenSelector(StrategyPriority)

	tokens := []*token.Token{
		{TokenID: "token1", Priority: 1},
		{TokenID: "token2", Priority: 3},
		{TokenID: "token3", Priority: 2},
	}

	selected := selector.SelectToken(tokens)
	if selected == nil {
		t.Fatal("Expected token, got nil")
	}

	if selected.TokenID != "token2" {
		t.Errorf("Expected highest priority token 'token2', got '%s'", selected.TokenID)
	}
}

func TestTokenSelector_EmptyTokens(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)

	selected := selector.SelectToken([]*token.Token{})
	if selected != nil {
		t.Error("Expected nil for empty tokens")
	}
}

func TestTokenSelector_NilTokens(t *testing.T) {
	selector := NewTokenSelector(StrategyRoundRobin)

	selected := selector.SelectToken(nil)
	if selected != nil {
		t.Error("Expected nil for nil tokens")
	}
}

func TestTokenSelector_DefaultStrategy(t *testing.T) {
	selector := NewTokenSelector("unknown")

	tokens := []*token.Token{
		{TokenID: "token1", Priority: 1},
		{TokenID: "token2", Priority: 1},
	}

	selected := selector.SelectToken(tokens)
	if selected == nil {
		t.Fatal("Expected token, got nil")
	}
}
