package golanglibs

import (
	"testing"
)

func TestRandomChoice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	for range Range(10) {
		Print(Random.Choice(s))
	}
}
