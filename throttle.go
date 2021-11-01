package throttle

import (
	"context"
	"errors"
)

type ctxType int

const CtxKey = ctxType(0x0)

func SetContext(ctx context.Context, identifier string) context.Context {
	return context.WithValue(ctx, CtxKey, identifier)
}

func FromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(CtxKey).(string)
	return s, ok
}

type KeyGetter interface {
	FromContext(ctx context.Context) string
}

type Throttle struct {
	l Limiter
}

func NewThrottle(l Limiter) Throttle {
	return Throttle{l: l}
}

func (t *Throttle) HandleContext(ctx context.Context) error {
	key, ok := FromContext(ctx)
	if !ok {
		return errors.New("missing key from header")
	}

	errc := make(chan error, 1)

	go func() {
		allowed, err := t.l.IsAllowed(ctx, key)
		if err != nil {
			errc <- err
		}

		if !allowed {
			errc <- errors.New("limit exceeded")
		}

		errc <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-errc:
		return e
	}
}
