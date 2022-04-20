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

func TestIn(t *testing.T) {
	Print(String("a").In("abc"))         // true
	Print(String("a").In("def"))         // false
	Print(String("a").In(String("abc"))) // true
	Print(String("a").In(String("def"))) // false
}

func TestDetectLang(t *testing.T) {
	Print(String("twitter is the best").Language())
}

func TestSimilar(t *testing.T) {
	s1 := `aio新手体验，
	A8哈希II超级活动，✈ @A8A8A88888
		   代理高扶持II最新 II牛牛II即将上线
			  尽情期待II我们A8客户II永久享受II我们的II高级待遇
	   详情II咨询上面客服II  ✈ @A8hash8868`
	s2 := `oaj新手体验，
	A8哈希II超级活动，✈ @A8A8A88888
		   代理高扶持II最新 II牛牛II即将上线
			  尽情期待II我们A8客户II永久享受II我们的II高级待遇
	   详情II咨询上面客服II  ✈ @A8hash8868`
	Print(String(s1).Similar(s2))
}
