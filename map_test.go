package golanglibs

import "testing"

func TestMap(t *testing.T) {
	m := map[string]string{
		"a": "b",
		"c": "d",
	}
	Print(Map(m).Has("a"))
	Print(Map(m).Keys())
}
