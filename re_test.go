package golanglibs

import (
	"testing"
)

func TestReFindall(t *testing.T) {
	Lg.Debug(Re.FindAll("[0-9]+", "a123b456"))
}
