package main

import (
	"math"
	"time"
)

// RetryPolicy decides whether a failed send should be retried, and
// how long to wait before the next attempt (exponential backoff).
type RetryPolicy struct {
	MaxRetries int
	BaseDelay  time.Duration
}

func NewRetryPolicy(maxRetry int, baseDelay time.Duration) *RetryPolicy {
	return &RetryPolicy{MaxRetries: maxRetry, BaseDelay: baseDelay}
}

// ShouldRetry returns true if we haven't exhausted our retry budget.
// In production this would also inspect the error type - e.g. don't
// retry on a 4xx "invalid phone number" but do retry on a timeout.
func (r *RetryPolicy) ShouldRetry(attempt int) bool {
	return attempt < r.MaxRetries
}

// NextDelay returns an exponential backoff duration: base * 2^attempt.
func (r *RetryPolicy) NextDelay(attempt int) time.Duration {
	factor := math.Pow(2, float64(attempt))
	return time.Duration(float64(r.BaseDelay) * factor)
}
