package ratelimiter

import "time"

// Rate represents a configs for the rate limiter.
type Rate struct {
	MaxRequest int           `json:"max_request,omitempty"`
	Duration   time.Duration `json:"duration,omitempty"`
}

// Stat holds a current state of each request.
type Stat struct {
	TotalHit int       `json:"total_hit,omitempty"`
	ResetAt  time.Time `json:"reset_at"`
}
