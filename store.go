package throttle

import "context"

type Store interface {
	SetStat(ctx context.Context, key string, stat Stat) error
	SetRate(ctx context.Context, key string, rate Rate) error
	GetRate(ctx context.Context, key string) (*Rate, error)
	GetStat(ctx context.Context, key string) (*Stat, error)
	Reset(ctx context.Context, key string) error
	Increment(ctx context.Context, key string) error
}
