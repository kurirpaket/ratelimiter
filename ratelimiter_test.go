package ratelimiter

import (
	"context"
	"errors"
	"testing"
)

func TestRateLimiter_HandleContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tableTest := []struct {
		desc string
		a    Allower
		ctx  context.Context
		err  error
	}{
		{
			desc: "should be canceled by deadline",
			a: AllowFunc(func(ctx context.Context, key string) error {
				return context.DeadlineExceeded
			}),
			ctx: SetContext(context.Background(), "example"),
			err: context.DeadlineExceeded,
		},
		{
			desc: "should be returns missing key",
			a: AllowFunc(func(ctx context.Context, key string) error {
				return nil
			}),
			ctx: context.Background(),
			err: ErrMissingKey,
		},
		{
			desc: "should be canceled by trigger",
			a: AllowFunc(func(ctx context.Context, key string) error {
				return nil
			}),
			ctx: SetContext(ctx, "example"),
			err: context.Canceled,
		},
	}

	for _, tt := range tableTest {
		tc := tt
		t.Run(tc.desc, func(t *testing.T) {
			rateLimiter := New(tc.a)
			if err := rateLimiter.HandleContext(tc.ctx); err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("expecting %v but got %v", tc.err, err)
				}
			}
		})
	}
}
