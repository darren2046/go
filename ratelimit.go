package golanglibs

import "go.uber.org/ratelimit"

type rateLimitStruct struct {
	rl ratelimit.Limiter
}

func getRateLimit(rate int) *rateLimitStruct {
	return &rateLimitStruct{
		rl: ratelimit.New(rate),
	}
}

func (m *rateLimitStruct) Take() {
	m.rl.Take()
}
