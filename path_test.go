package golanglibs

import "testing"

func TestPathJoin(t *testing.T) {
	Print(Os.Path.Join("/a/b/c", "/d/e/f"))
}
