package golanglibs

import (
	"testing"
)

func TestSub(t *testing.T) {
	if String("1234567890").Sub(0, 3).Get() != "123" {
		t.Error("Error while Sub")
	}
	if String("1234567890").Sub(2, 9).Get() != "3456789" {
		t.Error("Error while Sub")
	}
	if String("1234567890").Sub(7, 18).Get() != "890" {
		t.Error("Error while Sub")
	}
}

func TestChunk(t *testing.T) {
	s := String("1234567890").Chunk(3)
	for idx, str := range []string{"123", "456", "789", "0"} {
		if s[idx].S != str {
			t.Error("Error while Chunk")
		}
	}
}

func TestChr(t *testing.T) {
	if Chr(97) != "a" {
		t.Error("Error while Chr")
	}
}

func TestJoin(t *testing.T) {
	if String(".").Join([]string{"a", "b", "c"}).S != "a.b.c" {
		t.Error("Error while Join 1")
	}
	if String(".").Join(
		String("a,b,c").Split(","),
	).S != "a.b.c" {
		t.Error("Error while Join 2")
	}
}

func TestSplit(t *testing.T) {
	Print(String("a,b,c").Split(","))
}

func TestUtf8Split(t *testing.T) {
	Print(String("add陈四民w(*(f").Utf8Split())
}

func TestIsAscii(t *testing.T) {
	Print(String("add陈四民w(*(f").IsAscii())
	Print(String("addw(*(f").IsAscii())
}
