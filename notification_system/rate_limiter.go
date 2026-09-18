package main

import "sync"

// RateLimiter is a simple per-process token bucket. It protects
// downstream providers (SES/Twilio/FCM) from being overwhelmed.
// (In a real distributed deployment this would live in Redis so
// every instance shares the same bucket - see the HLD section.)
type RateLimiter struct {
	mu     sync.Mutex
	tokens int
	max    int
}

func NewRateLimiter(max int) *RateLimiter {
	return &RateLimiter{tokens: max, max: max}
}

// Allow consumes one token if available. Call Refill periodically
// (e.g. via time.Ticker in a background goroutine) to top up tokens.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.tokens <= 0 {
		return false
	}
	rl.tokens--
	return true
}

func (rl *RateLimiter) Refill() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.tokens = rl.max
}
