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
	StrategyCostFirst  TokenStrategy = "cost-first"
)

type TokenSelector struct {
	strategy    TokenStrategy
	costMatrix  map[string]CostConfig
	costTracker *CostTracker
	mu          sync.Mutex
	counter     uint64
}

func NewTokenSelector(strategy TokenStrategy) *TokenSelector {
	return &TokenSelector{
		strategy: strategy,
	}
}

func NewTokenSelectorWithCost(strategy TokenStrategy, costMatrix map[string]CostConfig) *TokenSelector {
	return &TokenSelector{
		strategy:   strategy,
		costMatrix: costMatrix,
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
	case StrategyCostFirst:
		return s.selectCostFirst(tokens)
	default:
		return s.selectRoundRobin(tokens)
	}
}

func (s *TokenSelector) selectRoundRobin(tokens []*token.Token) *token.Token {
	idx := atomic.AddUint64(&s.counter, 1)
	return tokens[(idx-1)%uint64(len(tokens))]
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

func (s *TokenSelector) selectCostFirst(tokens []*token.Token) *token.Token {
	if s.costMatrix == nil || len(tokens) == 0 {
		return tokens[0]
	}

	lowestCost := float64(0)
	selected := tokens[0]
	first := true

	for _, t := range tokens {
		var modelCost float64
		if len(t.AllowedModels) > 0 {
			modelCost = s.getModelCost(t.AllowedModels[0])
		}
		if first || modelCost < lowestCost {
			lowestCost = modelCost
			selected = t
			first = false
		}
	}

	return selected
}

func (s *TokenSelector) getModelCost(model string) float64 {
	if s.costMatrix == nil {
		return 0
	}
	cfg, ok := s.costMatrix[model]
	if !ok {
		return 0
	}
	return cfg.Input + cfg.Output
}
