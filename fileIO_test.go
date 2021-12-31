package golanglibs

import "testing"

func TestFileWrite(t *testing.T) {
	fd := Open("test.lllll", "w")
	fd.Write(String("stringStruct\n"))
	fd.Write([]byte("byte\n"))
	fd.Write("string")
	fd.Close()
}
