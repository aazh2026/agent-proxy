package routing

import (
	"sync"
	"time"
)

type QuotaType string

const (
	QuotaRequests QuotaType = "requests"
	QuotaTokens   QuotaType = "tokens"
	QuotaCost     QuotaType = "cost"
)

type QuotaPeriod string

const (
	PeriodDaily   QuotaPeriod = "daily"
	PeriodMonthly QuotaPeriod = "monthly"
)

type QuotaConfig struct {
	Type   QuotaType   `json:"type"`
	Period QuotaPeriod `json:"period"`
	Limit  int64       `json:"limit"`
	WarnAt float64     `json:"warn_at"`
}

type QuotaTracker struct {
	mu       sync.RWMutex
	counters map[string]*QuotaCounter
	configs  map[string]*QuotaConfig
}

type QuotaCounter struct {
	Requests int64
	Tokens   int64
	Cost     float64
	ResetAt  time.Time
}

func NewQuotaTracker() *QuotaTracker {
	return &QuotaTracker{
		counters: make(map[string]*QuotaCounter),
		configs:  make(map[string]*QuotaConfig),
	}
}

func (t *QuotaTracker) SetQuota(userID string, config *QuotaConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.configs[userID] = config
}

func (t *QuotaTracker) getOrCreateCounter(userID string) *QuotaCounter {
	counter, ok := t.counters[userID]
	if !ok {
		counter = &QuotaCounter{
			ResetAt: t.calculateResetTime(PeriodDaily),
		}
		t.counters[userID] = counter
	}

	if time.Now().After(counter.ResetAt) {
		counter.Requests = 0
		counter.Tokens = 0
		counter.Cost = 0
		counter.ResetAt = t.calculateResetTime(PeriodDaily)
	}

	return counter
}

func (t *QuotaTracker) GetCounter(userID string) *QuotaCounter {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.getOrCreateCounter(userID)
}

func (t *QuotaTracker) IncrementRequests(userID string, count int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	counter := t.getOrCreateCounter(userID)
	config := t.configs[userID]

	if config != nil && config.Type == QuotaRequests {
		if counter.Requests+count > config.Limit {
			return ErrQuotaExceeded
		}
	}

	counter.Requests += count
	return nil
}

func (t *QuotaTracker) IncrementTokens(userID string, tokens int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	counter := t.GetCounter(userID)
	config := t.configs[userID]

	if config != nil && config.Type == QuotaTokens {
		if counter.Tokens+tokens > config.Limit {
			return ErrQuotaExceeded
		}
	}

	counter.Tokens += tokens
	return nil
}

func (t *QuotaTracker) IncrementCost(userID string, cost float64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	counter := t.GetCounter(userID)
	config := t.configs[userID]

	if config != nil && config.Type == QuotaCost {
		if counter.Cost+cost > float64(config.Limit) {
			return ErrQuotaExceeded
		}
	}

	counter.Cost += cost
	return nil
}

func (t *QuotaTracker) CheckQuota(userID string) error {
	t.mu.RLock()
	config := t.configs[userID]
	t.mu.RUnlock()

	if config == nil {
		return nil
	}

	// Safely get counter with exclusive lock to avoid recursive lock acquisition.
	counter := t.GetCounter(userID)

	switch config.Type {
	case QuotaRequests:
		if counter.Requests >= config.Limit {
			return ErrQuotaExceeded
		}
	case QuotaTokens:
		if counter.Tokens >= config.Limit {
			return ErrQuotaExceeded
		}
	case QuotaCost:
		if counter.Cost >= float64(config.Limit) {
			return ErrQuotaExceeded
		}
	}

	return nil
}

func (t *QuotaTracker) calculateResetTime(period QuotaPeriod) time.Time {
	now := time.Now()
	switch period {
	case PeriodDaily:
		return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	case PeriodMonthly:
		return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	default:
		return now.Add(24 * time.Hour)
	}
}
