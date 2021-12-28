package golanglibs

import (
	"fmt"
	"testing"
)

func TestFakeNameChinese(t *testing.T) {
	for range Range(50) {
		fmt.Println(Funcs.FakeNameChinese())
	}
}

func TestFakeNameEnglish(t *testing.T) {
	for range Range(50) {
		fmt.Println(Funcs.FakeNameEnglish())
	}
}
