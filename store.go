package ratelimiter

import "context"

// RateStatCounter is an interface for Rate and Stat storage and also keep track
// the stat counter.
type RateStatCounter interface {
	RateSetGetter
	StatSetGetter
	Counter
}

// Counter knows how to manage your rate counter.
type Counter interface {
	// Increment increments the stat counter.
	Increment(ctx context.Context, key string) error

	// Reset resets the state counter into initial stat.
	Reset(ctx context.Context, key string) error
}

// RateSetGetter is an interface for get and setting rate.
type RateSetGetter interface {
	// SetRate updates a current rate.
	SetRate(ctx context.Context, key string, rate *Rate) error
	// Rate returns the current rate.
	Rate(ctx context.Context, key string) (*Rate, error)
}

// StatSetGetter is an interface for get and setting stat.
type StatSetGetter interface {
	// SetStat updates a current stat.
	SetStat(ctx context.Context, key string, stat *Stat) error

	// Stat returns the current stat.
	Stat(ctx context.Context, key string) (*Stat, error)
}
