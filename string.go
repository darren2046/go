package golanglibs

import (
	"bytes"
	"math/rand"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/denisbrodbeck/striphtmltags"
)

type stringStruct struct {
	s string
}

func String(s string) *stringStruct {
	return &stringStruct{s: s}
}

// Return the final string
func (s *stringStruct) Get() string {
	return s.s
}

// Return the substring of the string
func (s *stringStruct) Sub(start, end int) *stringStruct {
	s.s = s.sub(start, end)
	return s
}

func (s *stringStruct) sub(start, end int) string {
	start_str_idx := 0
	i := 0
	for j := range s.s {
		if i == start {
			start_str_idx = j
		}
		if i == end {
			return s.s[start_str_idx:j]
		}
		i++
	}
	return s.s[start_str_idx:]
}

// Return the length of the string
func (s *stringStruct) Len() int {
	return len(s.s)
}

// Reverse the string
func (s *stringStruct) Reverse() string {
	runes := []rune(s.s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Split the string with specific length
func (s *stringStruct) Chunk(length int) (res []string) {
	start := 0
	for s.Len() > start {
		substr := s.sub(start, start+length)
		res = append(res, substr)
		start = start + length
	}
	return
}

// Calculate the length of a UTF-8 string
func (s *stringStruct) Utf8Len() int {
	return utf8.RuneCountInString(s.s)
}

// Repeat a string
func (s *stringStruct) Repeat(count int) *stringStruct {
	s.s = strings.Repeat(s.s, count)
	return s
}

// Shuffle a string
func (s *stringStruct) Shuffle() *stringStruct {
	runes := []rune(s.s)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ss := make([]rune, len(runes))
	for i, v := range r.Perm(len(runes)) {
		ss[i] = runes[v]
	}
	s.s = string(ss)
	return s
}

func (s *stringStruct) Index(substr string) int {
	return strings.Index(s.s, substr)
}

func (s *stringStruct) Replace(old, new string) *stringStruct {
	s.s = strings.ReplaceAll(s.s, old, new)
	return s
}

func (s *stringStruct) Upper() *stringStruct {
	s.s = strings.ToUpper(s.s)
	return s
}

func (s *stringStruct) Lower() *stringStruct {
	s.s = strings.ToLower(s.s)
	return s
}

func (s *stringStruct) Join(pieces []string) *stringStruct {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(s.s)
		}
	}
	s.s = buf.String()
	return s
}

func (s *stringStruct) Strip(characterMask ...string) *stringStruct {
	if len(characterMask) == 0 {
		s.s = strings.TrimSpace(s.s)
	}
	s.s = strings.Trim(s.s, characterMask[0])
	return s
}

func (s *stringStruct) Split(sep ...string) []string {
	var a []string
	if len(sep) != 0 {
		for _, v := range strings.Split(s.s, sep[0]) {
			a = append(a, String(v).Strip().Get())
		}
	} else {
		for _, v := range strings.Split(s.s, " ") {
			if String(v).Strip().Get() != "" {
				a = append(a, String(v).Strip().Get())
			}
		}
	}

	return a
}

func (s *stringStruct) Count(substr string) int {
	return strings.Count(s.s, substr)
}

func (s *stringStruct) EndsWith(substr string) (res bool) {
	if len(substr) <= len(s.s) && s.s[len(s.s)-len(substr):] == substr {
		res = true
	} else {
		res = false
	}
	return
}

func (s *stringStruct) StartsWith(substr string) (res bool) {
	if len(substr) <= len(s.s) && s.s[:len(substr)] == substr {
		res = true
	} else {
		res = false
	}
	return
}

func (s *stringStruct) Splitlines(strip ...bool) []string {
	var a []string
	tostrip := false
	if len(strip) != 0 {
		tostrip = strip[0]
	}
	for _, v := range strings.Split(s.s, "\n") {
		if tostrip {
			a = append(a, String(v).Strip().Get())
		} else {
			a = append(a, String(v).RStrip("\r").Get())
		}
	}

	return a
}

func (s *stringStruct) In(str string) bool {
	if String(str).Index(s.s) != -1 {
		return true
	}
	return false
}

func (s *stringStruct) LStrip(characterMask ...string) *stringStruct {
	if len(characterMask) == 0 {
		s.s = strings.TrimLeftFunc(s.s, unicode.IsSpace)
	}
	s.s = strings.TrimLeft(s.s, characterMask[0])
	return s
}

func (s *stringStruct) RStrip(characterMask ...string) *stringStruct {
	if len(characterMask) == 0 {
		s.s = strings.TrimRightFunc(s.s, unicode.IsSpace)
	}
	s.s = strings.TrimRight(s.s, characterMask[0])
	return s
}

func (ss *stringStruct) Isdigit() bool {
	str := ss.s
	if str == "" {
		return false
	}
	// Trim any whitespace
	str = strings.TrimSpace(str)
	if str[0] == '-' || str[0] == '+' {
		if len(str) == 1 {
			return false
		}
		str = str[1:]
	}
	// hex
	if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
		for _, h := range str[2:] {
			if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
				return false
			}
		}
		return true
	}
	// 0-9, Point, Scientific
	p, s, l := 0, 0, len(str)
	for i, v := range str {
		if v == '.' { // Point
			if p > 0 || s > 0 || i+1 == l {
				return false
			}
			p = i
		} else if v == 'e' || v == 'E' { // Scientific
			if i == 0 || s > 0 || i+1 == l {
				return false
			}
			s = i
		} else if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

func (s *stringStruct) hasChinese() bool {
	var count int
	for _, v := range s.s {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

func (s *stringStruct) Filter(charts ...string) *stringStruct {
	var res string
	strb := []byte(s.s)
	var chartsb []byte
	if len(charts) != 0 {
		chartsb = []byte(charts[0])
	} else {
		chartsb = []byte("1234567890_qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM.")
	}
	for _, c := range strb {
		if Array(chartsb).Has(c) {
			res = res + string(c)
		}
	}
	s.s = res
	return s
}

func (s *stringStruct) RemoveHtmlTag() *stringStruct {
	s.s = striphtmltags.StripTags(s.s)
	return s
}

func (s *stringStruct) RemoveNonUTF8Character() *stringStruct {
	if !utf8.ValidString(s.s) {
		v := make([]rune, 0, len(s.s))
		for i, r := range s.s {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s.s[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		s.s = string(v)
	}
	return s
}
