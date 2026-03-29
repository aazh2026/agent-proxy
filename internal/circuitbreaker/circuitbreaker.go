package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

type CircuitBreaker struct {
	mu              sync.RWMutex
	state           State
	failures        int
	successes       int
	lastFailure     time.Time
	config          *Config
	halfOpenAllowed bool
}

type Config struct {
	FailureThreshold int
	SuccessThreshold int
	Timeout          time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		FailureThreshold: 5,
		SuccessThreshold: 3,
		Timeout:          30 * time.Second,
	}
}

func New(config *Config) *CircuitBreaker {
	if config == nil {
		config = DefaultConfig()
	}
	return &CircuitBreaker{
		state:  StateClosed,
		config: config,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.Allow() {
		return ErrCircuitOpen
	}

	err := fn()
	cb.Record(err)
	return err
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailure) > cb.config.Timeout {
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = StateHalfOpen
			cb.halfOpenAllowed = true
			cb.mu.Unlock()
			cb.mu.RLock()
			return true
		}
		return false
	case StateHalfOpen:
		return cb.halfOpenAllowed
	default:
		return false
	}
}

func (cb *CircuitBreaker) Record(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.recordFailure()
	} else {
		cb.recordSuccess()
	}
}

func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	cb.lastFailure = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.config.FailureThreshold {
			cb.state = StateOpen
		}
	case StateHalfOpen:
		cb.state = StateOpen
		cb.halfOpenAllowed = false
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.successes++

	switch cb.state {
	case StateHalfOpen:
		if cb.successes >= cb.config.SuccessThreshold {
			cb.state = StateClosed
			cb.failures = 0
			cb.successes = 0
			cb.halfOpenAllowed = false
		}
	case StateClosed:
		cb.failures = 0
	}
}

func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.halfOpenAllowed = false
}

type Metrics struct {
	TotalRequests  int64
	SuccessfulReqs int64
	FailedRequests int64
	RejectedReqs   int64
	StateChanges   int64
}

func (cb *CircuitBreaker) Metrics() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"state":        cb.state.String(),
		"failures":     cb.failures,
		"successes":    cb.successes,
		"last_failure": cb.lastFailure,
	}
}
