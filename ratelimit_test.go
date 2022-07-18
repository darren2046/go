package golanglibs

import (
	"fmt"
	"testing"
)

func TestRateLimit(t *testing.T) {
	rl := getRateLimit(3)

	for _, i := range Range(10) {
		fmt.Print(i)
		rl.Take()
	}
}
