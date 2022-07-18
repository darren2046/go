package golanglibs

import "go.uber.org/ratelimit"

type RateLimitStruct struct {
	rl    ratelimit.Limiter
	limit int
}

// It creates a new rate limiter with a rate of `rate` and returns a pointer to a `RateLimitStruct`
// If the rate is 0 then no limit
func getRateLimit(rate int) *RateLimitStruct {
	return &RateLimitStruct{
		rl:    ratelimit.New(rate),
		limit: rate,
	}
}

func (m *RateLimitStruct) Take() {
	if m.limit != 0 {
		m.rl.Take()
	}
}
