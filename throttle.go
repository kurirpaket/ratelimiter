package throttle

import (
	"context"
	"errors"
)

var (
	// ErrMissingKey is an error that occurred when key is missing from context.
	ErrMissingKey = errors.New("throttle: missing key from header")
	// ErrLimitExceeded is an error that occurred when rate limit is exceeded.
	ErrLimitExceeded = errors.New("throttle: limit exceeded")
)

// ctxType is a type for internal context.
type ctxType int

// CtxKey is key that can be used for storing stat/rate into context.
const CtxKey = ctxType(0x0)

// SetContext sets the given context with identifier value.
func SetContext(ctx context.Context, identifier string) context.Context {
	return context.WithValue(ctx, CtxKey, identifier)
}

// FromContext gets the identifier value from context.
func FromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(CtxKey).(string)
	return s, ok
}

// Throttle is throttle.
type Throttle struct {
	a Allower
}

// New creates a new Throttle.
func New(a Allower) *Throttle {
	return &Throttle{a: a}
}

// HandleContext checks if the given context is allowed.
func (t *Throttle) HandleContext(ctx context.Context) error {
	key, ok := FromContext(ctx)
	if !ok {
		return ErrMissingKey
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- t.a.Allow(ctx, key)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-errChan:
		return e
	}
}
