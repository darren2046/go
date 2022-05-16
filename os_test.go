package golanglibs

import "testing"

func TestBasename(t *testing.T) {
	Os.Path.Basename("/path/to/dir")
}

func TestChmod(t *testing.T) {
	Open("testfile", "w").Write("test").Close()
	Os.Chmod("testfile", 755)
}
