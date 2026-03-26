package routing

import (
	"sync"
	"sync/atomic"

	"github.com/openclaw/agent-proxy/internal/token"
)

type TokenStrategy string

const (
	StrategyRoundRobin TokenStrategy = "round-robin"
	StrategyWeighted   TokenStrategy = "weighted"
	StrategyPriority   TokenStrategy = "priority"
)

type TokenSelector struct {
	strategy TokenStrategy
	mu       sync.Mutex
	counter  uint64
}

func NewTokenSelector(strategy TokenStrategy) *TokenSelector {
	return &TokenSelector{
		strategy: strategy,
	}
}

func (s *TokenSelector) SelectToken(tokens []*token.Token) *token.Token {
	if len(tokens) == 0 {
		return nil
	}

	switch s.strategy {
	case StrategyRoundRobin:
		return s.selectRoundRobin(tokens)
	case StrategyWeighted:
		return s.selectWeighted(tokens)
	case StrategyPriority:
		return s.selectPriority(tokens)
	default:
		return s.selectRoundRobin(tokens)
	}
}

func (s *TokenSelector) selectRoundRobin(tokens []*token.Token) *token.Token {
	idx := atomic.AddUint64(&s.counter, 1)
	return tokens[idx%uint64(len(tokens))]
}

func (s *TokenSelector) selectWeighted(tokens []*token.Token) *token.Token {
	totalWeight := 0
	for _, t := range tokens {
		totalWeight += t.Priority
	}

	if totalWeight == 0 {
		return s.selectRoundRobin(tokens)
	}

	idx := atomic.AddUint64(&s.counter, 1)
	target := int(idx % uint64(totalWeight))

	current := 0
	for _, t := range tokens {
		current += t.Priority
		if current > target {
			return t
		}
	}

	return tokens[0]
}

func (s *TokenSelector) selectPriority(tokens []*token.Token) *token.Token {
	if len(tokens) == 0 {
		return nil
	}

	selected := tokens[0]
	for _, t := range tokens[1:] {
		if t.Priority > selected.Priority {
			selected = t
		}
	}
	return selected
}
