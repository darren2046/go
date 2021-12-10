package golanglibs

import "testing"

func TestInArray(t *testing.T) {
	if Array([]string{"123", "abc", "def"}).Has("abc") != true {
		t.Error("Error while check InArray")
	}
}
