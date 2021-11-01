package throttle

import (
	"context"
	"fmt"
	"time"
)

const DefaultKey = "default"

// Allower is a contract that needed to implement before using the Throttle.
type Allower interface {
	// Allow returns nil if the given context is allowed.
	Allow(ctx context.Context, key string) error
}

// AllowFunc is an adapter for creating Allower from function.
type AllowFunc func(ctx context.Context, key string) error

// Allow implements the Allower interface.
func (f AllowFunc) Allow(ctx context.Context, key string) error {
	return f(ctx, key)
}

// limiter knows how to limit your process.
type limiter struct {
	rate    *Rate
	storage RateStatCounter
}

// NewRateLimiter creates a new rate limiter that satisfied the Allower interface.
func NewRateLimiter(defaultRate *Rate, storage RateStatCounter) *limiter {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rate, err := storage.Rate(ctx, DefaultKey)
	if err != nil {
		_ = storage.SetRate(ctx, DefaultKey, defaultRate)
		rate = defaultRate
	}

	return &limiter{
		rate:    rate,
		storage: storage,
	}
}

// Allow implements the Allower interface.
func (l *limiter) Allow(ctx context.Context, key string) error {
	stat, err := l.storage.Stat(ctx, key)
	if err != nil {
		stat := Stat{
			TotalHit: 1,
			ResetAt:  time.Now(),
		}

		if err := l.storage.SetStat(ctx, key, &stat); err != nil {
			return fmt.Errorf("%w: setting stat", err)
		}

		return nil
	}

	rate, err := l.storage.Rate(ctx, key)
	if err != nil {
		rate = l.rate
	}

	if time.Since(stat.ResetAt) > rate.Duration {
		if err := l.storage.Reset(ctx, key); err != nil {
			return fmt.Errorf("%w: reseting counter", err)
		}

		return nil
	}

	if stat.TotalHit+1 > rate.MaxRequest {
		return ErrLimitExceeded
	}

	if err := l.storage.Increment(ctx, key); err != nil {
		return fmt.Errorf("%w: incrementing counter", err)
	}

	return nil
}
