package golanglibs

import (
	"testing"
)

func TestCmdWhich(t *testing.T) {
	Print(Cmd.Which("ls"))
	Try(func() {
		Print(Cmd.Which("NotExists"))
	}).Except(func(e error) {
		Print(e)
	})
}
