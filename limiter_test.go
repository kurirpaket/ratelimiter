package throttle

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mock struct {
	SetRateFunc   func(ctx context.Context, key string, rate *Rate) error
	RateFunc      func(ctx context.Context, key string) (*Rate, error)
	SetStatFunc   func(ctx context.Context, key string, stat *Stat) error
	StatFunc      func(ctx context.Context, key string) (*Stat, error)
	IncrementFunc func(ctx context.Context, key string) error
	ResetFunc     func(ctx context.Context, key string) error
}

func (m *mock) SetRate(ctx context.Context, key string, rate *Rate) error {
	return m.SetRateFunc(ctx, key, rate)
}

func (m *mock) Rate(ctx context.Context, key string) (*Rate, error) {
	return m.RateFunc(ctx, key)
}

func (m *mock) SetStat(ctx context.Context, key string, stat *Stat) error {
	return m.SetStatFunc(ctx, key, stat)
}

func (m *mock) Stat(ctx context.Context, key string) (*Stat, error) {
	return m.StatFunc(ctx, key)
}

func (m *mock) Increment(ctx context.Context, key string) error {
	return m.IncrementFunc(ctx, key)
}

func (m *mock) Reset(ctx context.Context, key string) error {
	return m.ResetFunc(ctx, key)
}

func TestNewRateLimiter(t *testing.T) {
	defaultRate := &Rate{}
	NewRateLimiter(defaultRate, &mock{
		SetRateFunc: func(ctx context.Context, key string, rate *Rate) error {
			if rate != defaultRate {
				t.Fatal("expecting using default rate")
			}
			return nil
		},
		RateFunc: func(ctx context.Context, key string) (*Rate, error) {
			return nil, errors.New("")
		},
		SetStatFunc:   nil,
		StatFunc:      nil,
		IncrementFunc: nil,
		ResetFunc:     nil,
	})
}

func TestLimiter_Allow(t *testing.T) {
	t.Run("expecting reset state", func(t *testing.T) {
		defaultRate := &Rate{}
		limiter := NewRateLimiter(defaultRate, &mock{
			SetRateFunc: func(ctx context.Context, key string, rate *Rate) error {
				t.Errorf("expecting not called")
				return nil
			},
			RateFunc: func(ctx context.Context, key string) (*Rate, error) {
				rate := &Rate{
					MaxRequest: 2,
					Duration:   time.Minute,
				}

				return rate, nil
			},
			SetStatFunc: func(ctx context.Context, key string, stat *Stat) error {
				return nil
			},
			StatFunc: func(ctx context.Context, key string) (*Stat, error) {
				return nil, errors.New("not exist")
			},
			IncrementFunc: func(ctx context.Context, key string) error {
				return nil
			},
			ResetFunc: func(ctx context.Context, key string) error {
				return nil
			},
		})

		if err := limiter.Allow(SetContext(context.Background(), "example"), "example"); err != nil {
			if !errors.Is(err, ErrLimitExceeded) {
				t.Errorf("expecting limit exceeded")
			}
		}

	})
	t.Run("expecting reset state but got an error", func(t *testing.T) {
		defaultRate := &Rate{}
		limiter := NewRateLimiter(defaultRate, &mock{
			SetRateFunc: func(ctx context.Context, key string, rate *Rate) error {
				t.Errorf("expecting not called")
				return nil
			},
			RateFunc: func(ctx context.Context, key string) (*Rate, error) {
				rate := &Rate{
					MaxRequest: 2,
					Duration:   time.Minute,
				}

				return rate, nil
			},
			SetStatFunc: func(ctx context.Context, key string, stat *Stat) error {
				return errors.New("failed")
			},
			StatFunc: func(ctx context.Context, key string) (*Stat, error) {
				return nil, errors.New("not exist")
			},
			IncrementFunc: func(ctx context.Context, key string) error {
				return nil
			},
			ResetFunc: func(ctx context.Context, key string) error {
				return nil
			},
		})

		if err := limiter.Allow(SetContext(context.Background(), "example"), "example"); err == nil {
			t.Errorf("expecting non nil error")
		}

	})
	t.Run("expecting using default rate when error get rate from storage is error", func(t *testing.T) {
		defaultRate := &Rate{
			MaxRequest: 2,
			Duration:   time.Minute,
		}
		limiter := NewRateLimiter(defaultRate, &mock{
			SetRateFunc: func(ctx context.Context, key string, rate *Rate) error {
				return nil
			},
			RateFunc: func(ctx context.Context, key string) (*Rate, error) {
				return nil, errors.New("failed")
			},
			SetStatFunc: func(ctx context.Context, key string, stat *Stat) error {
				return nil
			},
			StatFunc: func(ctx context.Context, key string) (*Stat, error) {
				return &Stat{}, nil
			},
			IncrementFunc: func(ctx context.Context, key string) error {
				return nil
			},
			ResetFunc: func(ctx context.Context, key string) error {
				return nil
			},
		})

		if err := limiter.Allow(SetContext(context.Background(), "example"), "example"); err != nil {
			t.Errorf("expecting  nil error")
		}
	})
	t.Run("expecting limit exceeded", func(t *testing.T) {
		defaultRate := &Rate{}
		limiter := NewRateLimiter(defaultRate, &mock{
			SetRateFunc: func(ctx context.Context, key string, rate *Rate) error {
				t.Errorf("expecting not called")
				return nil
			},
			RateFunc: func(ctx context.Context, key string) (*Rate, error) {
				rate := &Rate{
					MaxRequest: 2,
					Duration:   time.Minute,
				}

				return rate, nil
			},
			SetStatFunc: nil,
			StatFunc: func(ctx context.Context, key string) (*Stat, error) {
				stat := &Stat{
					TotalHit: 2,
					ResetAt:  time.Now(),
				}

				return stat, nil
			},
			IncrementFunc: func(ctx context.Context, key string) error {
				return nil
			},
			ResetFunc: func(ctx context.Context, key string) error {
				return nil
			},
		})

		if err := limiter.Allow(SetContext(context.Background(), "example"), "example"); err != nil {
			if !errors.Is(err, ErrLimitExceeded) {
				t.Errorf("expecting limit exceeded")
			}
		}

	})
	t.Run("expecting error when resetting is failed", func(t *testing.T) {
		defaultRate := &Rate{}
		limiter := NewRateLimiter(defaultRate, &mock{
			SetRateFunc: func(ctx context.Context, key string, rate *Rate) error {
				t.Errorf("expecting not called")
				return nil
			},
			RateFunc: func(ctx context.Context, key string) (*Rate, error) {
				rate := &Rate{
					MaxRequest: 2,
					Duration:   time.Minute,
				}

				return rate, nil
			},
			SetStatFunc: nil,
			StatFunc: func(ctx context.Context, key string) (*Stat, error) {
				stat := &Stat{
					TotalHit: 2,
					ResetAt:  time.Time{},
				}

				return stat, nil
			},
			IncrementFunc: func(ctx context.Context, key string) error {
				return nil
			},
			ResetFunc: func(ctx context.Context, key string) error {
				return errors.New("failed")
			},
		})

		if err := limiter.Allow(SetContext(context.Background(), "example"), "example"); err == nil {
			t.Errorf("expecting error nil")
		}
	})
	t.Run("expecting error when incrementing is failed", func(t *testing.T) {
		defaultRate := &Rate{}
		limiter := NewRateLimiter(defaultRate, &mock{
			SetRateFunc: func(ctx context.Context, key string, rate *Rate) error {
				t.Errorf("expecting not called")
				return nil
			},
			RateFunc: func(ctx context.Context, key string) (*Rate, error) {
				rate := &Rate{
					MaxRequest: 2,
					Duration:   time.Minute,
				}

				return rate, nil
			},
			SetStatFunc: nil,
			StatFunc: func(ctx context.Context, key string) (*Stat, error) {
				stat := &Stat{
					TotalHit: 0,
					ResetAt:  time.Now(),
				}

				return stat, nil
			},
			IncrementFunc: func(ctx context.Context, key string) error {
				return errors.New("failed")
			},
			ResetFunc: func(ctx context.Context, key string) error {
				return nil
			},
		})

		if err := limiter.Allow(SetContext(context.Background(), "example"), "example"); err == nil {
			t.Errorf("expecting non error nil")
		}
	})
	t.Run("expecting ok", func(t *testing.T) {
		defaultRate := &Rate{}
		limiter := NewRateLimiter(defaultRate, &mock{
			SetRateFunc: func(ctx context.Context, key string, rate *Rate) error {
				t.Errorf("expecting not called")
				return nil
			},
			RateFunc: func(ctx context.Context, key string) (*Rate, error) {
				rate := &Rate{
					MaxRequest: 2,
					Duration:   time.Minute,
				}

				return rate, nil
			},
			SetStatFunc: nil,
			StatFunc: func(ctx context.Context, key string) (*Stat, error) {
				stat := &Stat{
					TotalHit: 0,
					ResetAt:  time.Now(),
				}

				return stat, nil
			},
			IncrementFunc: func(ctx context.Context, key string) error {
				return nil
			},
			ResetFunc: func(ctx context.Context, key string) error {
				return nil
			},
		})

		if err := limiter.Allow(SetContext(context.Background(), "example"), "example"); err != nil {
			t.Errorf("expecting error nil")
		}
	})
}
