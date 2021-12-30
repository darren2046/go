package golanglibs

import "testing"

func TestFileWrite(t *testing.T) {
	k := Open("file.go").Read()
	Print(Typeof(k))
}
