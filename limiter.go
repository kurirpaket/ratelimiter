package throttle

import (
	"context"
	"time"
)

type Limiter interface {
	IsAllowed(ctx context.Context, key string) (bool, error)
}

type rateLimiter struct {
	*Rate
	store Store
}

func NewRateLimiter(defaultRate Rate, store Store) Limiter {
	rate, err := store.GetRate(context.Background(), "default")
	if err != nil {
		go store.SetRate(context.Background(), "default", defaultRate)
		rate = &defaultRate
	}

	rl := rateLimiter{rate, store}

	return rl
}

func (r rateLimiter) IsAllowed(ctx context.Context, key string) (bool, error) {
	currentStat, err := r.store.GetStat(ctx, key)
	if err != nil {
		err = r.store.SetStat(ctx, key, Stat{
			TotalHit: 1,
			ResetAt:  time.Now(),
		})
		if err != nil {
			return false, err
		}

		return true, nil
	}

	rate, err := r.store.GetRate(ctx, key)
	if err != nil {
		rate = r.Rate
	}

	if time.Since(currentStat.ResetAt) > rate.Duration {
		_ = r.store.Reset(ctx, key)
		return true, nil
	}

	if currentStat.TotalHit+1 > rate.MaxRequest {
		return false, nil
	}

	err = r.store.Increment(ctx, key)
	if err != nil {
		return false, err
	}

	return true, nil
}
