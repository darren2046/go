package golanglibs

import "testing"

func TestBasename(t *testing.T) {
	Os.Path.Basename("/path/to/dir")
}

func TestChmod(t *testing.T) {
	Os.Chmod("os.go", 644)
}
