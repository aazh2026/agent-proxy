package routing

import (
	"context"
	"strings"

	"github.com/openclaw/agent-proxy/internal/token"
)

type CostStrategy string

const (
	CostStrategyQualityFirst CostStrategy = "quality-first"
	CostStrategyCostFirst    CostStrategy = "cost-first"
	CostStrategyLatencyFirst CostStrategy = "latency-first"
	CostStrategyBalanced     CostStrategy = "balanced"
)

type CostConfig struct {
	Input  float64
	Output float64
}

type CostSelector struct {
	costMatrix map[string]CostConfig
	strategy   CostStrategy
	userPrefs  map[string]string
}

func NewCostSelector(costMatrix map[string]CostConfig, strategy string, userPrefs map[string]string) *CostSelector {
	s := CostStrategy(CostStrategyQualityFirst)
	switch strategy {
	case "cost-first":
		s = CostStrategyCostFirst
	case "latency-first":
		s = CostStrategyLatencyFirst
	case "balanced":
		s = CostStrategyBalanced
	}
	return &CostSelector{
		costMatrix: costMatrix,
		strategy:   s,
		userPrefs:  userPrefs,
	}
}

func (c *CostSelector) GetStrategyForUser(userID string) CostStrategy {
	if c.userPrefs != nil {
		if pref, ok := c.userPrefs[userID]; ok {
			switch pref {
			case "cost-first":
				return CostStrategyCostFirst
			case "latency-first":
				return CostStrategyLatencyFirst
			case "balanced":
				return CostStrategyBalanced
			}
		}
	}
	return c.strategy
}

func (c *CostSelector) CalculateCost(model string, promptTokens, completionTokens int) float64 {
	if c.costMatrix == nil {
		return 0
	}
	cfg, ok := c.costMatrix[model]
	if !ok {
		return 0
	}
	inputCost := float64(promptTokens) / 1_000_000 * cfg.Input
	outputCost := float64(completionTokens) / 1_000_000 * cfg.Output
	return inputCost + outputCost
}

func (c *CostSelector) SelectModel(ctx context.Context, availableModels []string, strategy CostStrategy) string {
	if len(availableModels) == 0 {
		return ""
	}

	switch strategy {
	case CostStrategyCostFirst:
		return c.selectCheapest(availableModels)
	case CostStrategyQualityFirst:
		return c.selectMostExpensive(availableModels)
	case CostStrategyLatencyFirst:
		return availableModels[0]
	case CostStrategyBalanced:
		return c.selectBalanced(availableModels)
	default:
		return availableModels[0]
	}
}

func (c *CostSelector) selectCheapest(models []string) string {
	lowestCost := float64(0)
	cheapest := models[0]
	first := true

	for _, m := range models {
		if c.costMatrix == nil {
			return m
		}
		cfg, ok := c.costMatrix[m]
		if !ok {
			continue
		}
		totalCost := cfg.Input + cfg.Output
		if first || totalCost < lowestCost {
			lowestCost = totalCost
			cheapest = m
			first = false
		}
	}
	return cheapest
}

func (c *CostSelector) selectMostExpensive(models []string) string {
	highestCost := float64(0)
	expensive := models[0]
	first := true

	for _, m := range models {
		if c.costMatrix == nil {
			return m
		}
		cfg, ok := c.costMatrix[m]
		if !ok {
			continue
		}
		totalCost := cfg.Input + cfg.Output
		if first || totalCost > highestCost {
			highestCost = totalCost
			expensive = m
			first = false
		}
	}
	return expensive
}

func (c *CostSelector) selectBalanced(models []string) string {
	if len(models) <= 1 {
		return models[0]
	}
	mid := len(models) / 2
	return models[mid]
}

type CostTracker struct {
	userCosts map[string]*UserCost
}

type UserCost struct {
	TotalCost        float64
	PromptTokens     int
	CompletionTokens int
	RequestCount     int
}

func NewCostTracker() *CostTracker {
	return &CostTracker{
		userCosts: make(map[string]*UserCost),
	}
}

func (t *CostTracker) RecordCost(userID string, cost float64, promptTokens, completionTokens int) {
	if t.userCosts[userID] == nil {
		t.userCosts[userID] = &UserCost{}
	}
	t.userCosts[userID].TotalCost += cost
	t.userCosts[userID].PromptTokens += promptTokens
	t.userCosts[userID].CompletionTokens += completionTokens
	t.userCosts[userID].RequestCount++
}

func (t *CostTracker) GetUserCost(userID string) *UserCost {
	if t.userCosts[userID] == nil {
		return &UserCost{}
	}
	return t.userCosts[userID]
}

func (t *CostTracker) ResetUserCost(userID string) {
	delete(t.userCosts, userID)
}

func SelectTokenByCost(tokens []*token.Token, costMatrix map[string]CostConfig, strategy CostStrategy, requestedModel string) *token.Token {
	if len(tokens) == 0 {
		return nil
	}

	if requestedModel == "" {
		return tokens[0]
	}

	selector := &CostSelector{
		costMatrix: costMatrix,
		strategy:   strategy,
	}

	provider := ""
	for _, t := range tokens {
		if t.Provider == "openai" && (requestedModel == "gpt-4o" || requestedModel == "gpt-4o-mini" || requestedModel == "gpt-3.5-turbo" || requestedModel == "gpt-4") {
			provider = "openai"
			break
		}
		if t.Provider == "anthropic" && (strings.HasPrefix(requestedModel, "claude-")) {
			provider = "anthropic"
			break
		}
		if t.Provider == "google" && (strings.HasPrefix(requestedModel, "gemini-")) {
			provider = "google"
			break
		}
	}

	for _, t := range tokens {
		if t.Provider == provider {
			return t
		}
	}

	if strategy == CostStrategyCostFirst && costMatrix != nil {
		availableModels := make([]string, 0)
		seen := make(map[string]bool)
		for _, t := range tokens {
			modelName := t.Provider + "-" + t.AllowedModels[0]
			if !seen[modelName] && len(t.AllowedModels) > 0 {
				seen[modelName] = true
				availableModels = append(availableModels, modelName)
			}
		}
		if len(availableModels) > 0 {
			selectedModel := selector.SelectModel(context.Background(), availableModels, strategy)
			for _, t := range tokens {
				if len(t.AllowedModels) > 0 && t.AllowedModels[0] == strings.TrimPrefix(selectedModel, t.Provider+"-") {
					return t
				}
			}
		}
	}

	return tokens[0]
}
