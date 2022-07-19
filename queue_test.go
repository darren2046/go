package golanglibs

import (
	"testing"
)

func TestQueue(t *testing.T) {
	qb := getQueue("data")

	qn := qb.New()

	if qn.Size() != 0 {
		t.Error("Size not correct")
	}
	qn.Put("value1")
	if qn.Size() != 1 {
		t.Error("Size not correct")
	}

	if qn.Get() != "value1" {
		t.Error("value not correct")
	}
	if qn.Size() != 0 {
		t.Error("Size not correct")
	}

	qb.Close()
}
