package routing

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*TokenBucket
	config   *RateLimitConfig
}

type RateLimitConfig struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	BurstSize         int `json:"burst_size"`
}

type TokenBucket struct {
	tokens     int
	lastRefill time.Time
}

func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	if config == nil {
		config = &RateLimitConfig{
			RequestsPerMinute: 60,
			BurstSize:         10,
		}
	}
	return &RateLimiter{
		limiters: make(map[string]*TokenBucket),
		config:   config,
	}
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	bucket, ok := r.limiters[key]
	if !ok {
		bucket = &TokenBucket{
			tokens:     r.config.BurstSize,
			lastRefill: time.Now(),
		}
		r.limiters[key] = bucket
	}

	r.refill(bucket)

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}
	return false
}

func (r *RateLimiter) refill(bucket *TokenBucket) {
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(elapsed.Seconds() * float64(r.config.RequestsPerMinute) / 60.0)

	if tokensToAdd > 0 {
		bucket.tokens += tokensToAdd
		if bucket.tokens > r.config.BurstSize {
			bucket.tokens = r.config.BurstSize
		}
		bucket.lastRefill = now
	}
}

func (r *RateLimiter) GetRemaining(key string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bucket, ok := r.limiters[key]
	if !ok {
		return r.config.BurstSize
	}

	r.refill(bucket)
	return bucket.tokens
}

func (r *RateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.limiters, key)
}
