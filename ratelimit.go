package golanglibs

import "go.uber.org/ratelimit"

type RateLimitStruct struct {
	rl ratelimit.Limiter
}

func getRateLimit(rate int) *RateLimitStruct {
	return &RateLimitStruct{
		rl: ratelimit.New(rate),
	}
}

func (m *RateLimitStruct) Take() {
	m.rl.Take()
}
