package golanglibs

import (
	"fmt"
	"testing"
)

func TestCSV(t *testing.T) {
	w := Tools.CSV.Writer("test.csv", "w")
	w.SetHeaders([]string{"first", "second", "third"})
	w.Write(map[string]string{
		"first":  "aaaaa",
		"second": "bbbbb",
		"third":  "ccccc",
	})
	w.Close()

	fmt.Println("Now csv content is")
	Os.System("cat test.csv")

	w = Tools.CSV.Writer("test.csv", "a")
	w.Write(map[string]string{
		"first":  "11111",
		"second": "22222",
		"third":  "33333",
	})
	w.Close()

	fmt.Println("Now csv content is")
	Os.System("cat test.csv")

	r := Tools.CSV.Reader("test.csv")
	fmt.Println("Read a row:")
	Print(r.Read())
	r.Close()

	fmt.Println("Read all:")
	r = Tools.CSV.Reader("test.csv")
	for i := range r.Readrows() {
		Print(i)
	}

	Os.Unlink("test.csv")
}
