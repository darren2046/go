package golanglibs

import (
	"testing"
)

func TestQueue(t *testing.T) {
	qb := getQueue("data")

	qn := qb.New()

	Lg.Debug("Size:", qn.Size())
	qn.Put("value1")
	Lg.Debug("Size:", qn.Size())
	qb.Close()

	qq := getQueue("data")
	qm := qq.New()
	Lg.Debug("Size:", qm.Size())
	Lg.Debug(qm.Get())
	Lg.Debug("Size:", qm.Size())

	qq.Destroy()

	if Os.Path.Exists("data") {
		t.Error("Faild")
	}
}
