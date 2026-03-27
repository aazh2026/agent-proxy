package routing

import (
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	limiter := NewRateLimiter(&RateLimitConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
	})

	for i := 0; i < 10; i++ {
		if !limiter.Allow("user1") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	if limiter.Allow("user1") {
		t.Error("Request 11 should be rejected")
	}
}

func TestRateLimiterRefill(t *testing.T) {
	limiter := NewRateLimiter(&RateLimitConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
	})

	for i := 0; i < 10; i++ {
		limiter.Allow("user1")
	}

	if limiter.Allow("user1") {
		t.Error("Should be rate limited")
	}

	time.Sleep(1 * time.Second)

	if !limiter.Allow("user1") {
		t.Error("Should have refilled tokens")
	}
}

func TestRateLimiterMultipleUsers(t *testing.T) {
	limiter := NewRateLimiter(&RateLimitConfig{
		RequestsPerMinute: 60,
		BurstSize:         5,
	})

	for i := 0; i < 5; i++ {
		if !limiter.Allow("user1") {
			t.Errorf("User1 request %d should be allowed", i+1)
		}
	}

	for i := 0; i < 5; i++ {
		if !limiter.Allow("user2") {
			t.Errorf("User2 request %d should be allowed", i+1)
		}
	}

	if limiter.Allow("user1") {
		t.Error("User1 should be rate limited")
	}
	if limiter.Allow("user2") {
		t.Error("User2 should be rate limited")
	}
}
