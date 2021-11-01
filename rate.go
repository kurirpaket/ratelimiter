package throttle

import "time"

type Rate struct {
	MaxRequest int           `json:"max_request,omitempty"`
	Duration   time.Duration `json:"duration,omitempty"`
}

type Stat struct {
	TotalHit int       `json:"total_hit,omitempty"`
	ResetAt  time.Time `json:"reset_at"`
}
