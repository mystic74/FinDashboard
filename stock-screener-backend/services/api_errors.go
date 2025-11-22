package services

import (
	"errors"
	"fmt"
	"time"
)

// API Error Types
var (
	// ErrRateLimited indicates the API rate limit has been exceeded
	ErrRateLimited = errors.New("rate limit exceeded")

	// ErrQuotaExceeded indicates the daily/monthly quota has been exceeded
	ErrQuotaExceeded = errors.New("API quota exceeded")

	// ErrPremiumRequired indicates the endpoint requires a premium subscription
	ErrPremiumRequired = errors.New("premium subscription required")

	// ErrInvalidAPIKey indicates the API key is invalid or missing
	ErrInvalidAPIKey = errors.New("invalid or missing API key")

	// ErrProviderUnavailable indicates the provider is temporarily unavailable
	ErrProviderUnavailable = errors.New("data provider unavailable")

	// ErrSymbolNotFound indicates the requested symbol was not found
	ErrSymbolNotFound = errors.New("symbol not found")
)

// APIError provides detailed error information from data providers
type APIError struct {
	Provider    string        // e.g., "alpha_vantage", "fmp", "yahoo"
	Code        string        // Provider-specific error code
	Message     string        // Human-readable message
	RetryAfter  time.Duration // Suggested retry delay (for rate limits)
	Recoverable bool          // Whether the error is recoverable with retry
	Underlying  error         // Underlying error type
}

func (e *APIError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("[%s] %s (retry after %v)", e.Provider, e.Message, e.RetryAfter)
	}
	return fmt.Sprintf("[%s] %s", e.Provider, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Underlying
}

// Is implements error matching for errors.Is
func (e *APIError) Is(target error) bool {
	return errors.Is(e.Underlying, target)
}

// NewRateLimitError creates a rate limit error with retry information
func NewRateLimitError(provider string, retryAfter time.Duration) *APIError {
	return &APIError{
		Provider:    provider,
		Code:        "RATE_LIMITED",
		Message:     "API rate limit exceeded",
		RetryAfter:  retryAfter,
		Recoverable: true,
		Underlying:  ErrRateLimited,
	}
}

// NewQuotaExceededError creates a quota exceeded error
func NewQuotaExceededError(provider string, quotaType string) *APIError {
	return &APIError{
		Provider:    provider,
		Code:        "QUOTA_EXCEEDED",
		Message:     fmt.Sprintf("%s quota exceeded", quotaType),
		Recoverable: false, // Need to wait for quota reset
		Underlying:  ErrQuotaExceeded,
	}
}

// NewPremiumRequiredError creates an error for premium-only features
func NewPremiumRequiredError(provider string, feature string) *APIError {
	return &APIError{
		Provider:    provider,
		Code:        "PREMIUM_REQUIRED",
		Message:     fmt.Sprintf("'%s' requires premium subscription", feature),
		Recoverable: false,
		Underlying:  ErrPremiumRequired,
	}
}

// NewInvalidAPIKeyError creates an invalid API key error
func NewInvalidAPIKeyError(provider string) *APIError {
	return &APIError{
		Provider:    provider,
		Code:        "INVALID_API_KEY",
		Message:     "Invalid or missing API key",
		Recoverable: false,
		Underlying:  ErrInvalidAPIKey,
	}
}

// NewProviderUnavailableError creates a provider unavailable error
func NewProviderUnavailableError(provider string, reason string) *APIError {
	return &APIError{
		Provider:    provider,
		Code:        "UNAVAILABLE",
		Message:     fmt.Sprintf("Provider unavailable: %s", reason),
		RetryAfter:  30 * time.Second,
		Recoverable: true,
		Underlying:  ErrProviderUnavailable,
	}
}

// RateLimiter manages API rate limiting
type RateLimiter struct {
	provider      string
	callsPerMin   int
	callsPerDay   int
	minCallCount  int
	dayCallCount  int
	lastMinReset  time.Time
	lastDayReset  time.Time
	retryAfter    time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(provider string, callsPerMin, callsPerDay int) *RateLimiter {
	now := time.Now()
	return &RateLimiter{
		provider:     provider,
		callsPerMin:  callsPerMin,
		callsPerDay:  callsPerDay,
		lastMinReset: now,
		lastDayReset: now,
	}
}

// CanMakeCall checks if a call can be made without hitting rate limits
func (r *RateLimiter) CanMakeCall() bool {
	now := time.Now()

	// Check if we're in a forced retry-after period
	if now.Before(r.retryAfter) {
		return false
	}

	// Reset minute counter if a minute has passed
	if now.Sub(r.lastMinReset) >= time.Minute {
		r.minCallCount = 0
		r.lastMinReset = now
	}

	// Reset day counter if a day has passed
	if now.Sub(r.lastDayReset) >= 24*time.Hour {
		r.dayCallCount = 0
		r.lastDayReset = now
	}

	// Check limits
	if r.callsPerMin > 0 && r.minCallCount >= r.callsPerMin {
		return false
	}
	if r.callsPerDay > 0 && r.dayCallCount >= r.callsPerDay {
		return false
	}

	return true
}

// RecordCall records that an API call was made
func (r *RateLimiter) RecordCall() {
	r.minCallCount++
	r.dayCallCount++
}

// SetRetryAfter sets a forced wait period (from API response headers)
func (r *RateLimiter) SetRetryAfter(duration time.Duration) {
	r.retryAfter = time.Now().Add(duration)
}

// GetWaitTime returns how long to wait before the next call
func (r *RateLimiter) GetWaitTime() time.Duration {
	now := time.Now()

	// Check forced retry-after first
	if now.Before(r.retryAfter) {
		return r.retryAfter.Sub(now)
	}

	// Check minute limit
	if r.callsPerMin > 0 && r.minCallCount >= r.callsPerMin {
		return time.Minute - now.Sub(r.lastMinReset)
	}

	// Check day limit
	if r.callsPerDay > 0 && r.dayCallCount >= r.callsPerDay {
		return 24*time.Hour - now.Sub(r.lastDayReset)
	}

	return 0
}

// GetStatus returns the current rate limit status
func (r *RateLimiter) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"provider":        r.provider,
		"minuteUsed":      r.minCallCount,
		"minuteLimit":     r.callsPerMin,
		"dayUsed":         r.dayCallCount,
		"dayLimit":        r.callsPerDay,
		"canMakeCall":     r.CanMakeCall(),
		"waitTime":        r.GetWaitTime().String(),
	}
}

// ProviderStatus tracks the health status of a data provider
type ProviderStatus struct {
	Name           string
	Available      bool
	LastError      error
	LastErrorTime  time.Time
	ConsecutiveFails int
	LastSuccessTime time.Time
}

// IsHealthy returns true if the provider is considered healthy
func (p *ProviderStatus) IsHealthy() bool {
	if !p.Available {
		return false
	}
	// Consider unhealthy if 3+ consecutive failures in last 5 minutes
	if p.ConsecutiveFails >= 3 && time.Since(p.LastErrorTime) < 5*time.Minute {
		return false
	}
	return true
}

// RecordSuccess records a successful API call
func (p *ProviderStatus) RecordSuccess() {
	p.Available = true
	p.ConsecutiveFails = 0
	p.LastSuccessTime = time.Now()
}

// RecordFailure records a failed API call
func (p *ProviderStatus) RecordFailure(err error) {
	p.LastError = err
	p.LastErrorTime = time.Now()
	p.ConsecutiveFails++

	// Mark as unavailable after too many failures
	if p.ConsecutiveFails >= 5 {
		p.Available = false
	}
}
