package golanglibs

import (
	"bytes"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/abadojack/whatlanggo"
	"github.com/denisbrodbeck/striphtmltags"
	"golang.org/x/exp/utf8string"
)

type stringStruct struct {
	S string
}

func String(s interface{}) *stringStruct {
	return &stringStruct{S: Str(s)}
}

// Return the final string
func (s *stringStruct) Get() string {
	return s.S
}

// Return the substring of the string
func (s *stringStruct) Sub(start, end int) *stringStruct {
	s.S = s.sub(start, end)
	return s
}

func (s *stringStruct) Has(substr string) bool {
	return strings.Contains(s.S, substr)
}

func (s *stringStruct) sub(start, end int) string {
	startStrIdx := 0
	i := 0
	for j := range s.S {
		if i == start {
			startStrIdx = j
		}
		if i == end {
			return s.S[startStrIdx:j]
		}
		i++
	}
	return s.S[startStrIdx:]
}

// Return the length of the string
func (s *stringStruct) Len() int {
	return len(s.S)
}

// Reverse the string
func (s *stringStruct) Reverse() *stringStruct {
	runes := []rune(s.S)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	s.S = string(runes)
	return s
}

// Split the string with specific length
func (s *stringStruct) Chunk(length int) (res []*stringStruct) {
	start := 0
	for s.Len() > start {
		substr := s.Sub(start, start+length)
		res = append(res, substr)
		start = start + length
	}
	return
}

// Calculate the length of a UTF-8 string
func (s *stringStruct) Utf8Len() int {
	return utf8.RuneCountInString(s.S)
}

// Repeat a string
func (s *stringStruct) Repeat(count int) *stringStruct {
	s.S = strings.Repeat(s.S, count)
	return s
}

// Shuffle a string
func (s *stringStruct) Shuffle() *stringStruct {
	runes := []rune(s.S)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ss := make([]rune, len(runes))
	for i, v := range r.Perm(len(runes)) {
		ss[i] = runes[v]
	}
	s.S = string(ss)
	return s
}

func (s *stringStruct) Index(substr string) int {
	return strings.Index(s.S, substr)
}

func (s *stringStruct) Replace(oldText, newText interface{}) *stringStruct {
	s.S = strings.ReplaceAll(s.S, Str(oldText), Str(newText))
	return s
}

func (s *stringStruct) Upper() *stringStruct {
	s.S = strings.ToUpper(s.S)
	return s
}

func (s *stringStruct) Lower() *stringStruct {
	s.S = strings.ToLower(s.S)
	return s
}

// pieces can be []string or []*stringStruct
func (s *stringStruct) Join(pieces interface{}) *stringStruct {
	var buf bytes.Buffer
	arr := reflect.ValueOf(pieces)
	l := arr.Len()

	for i := 0; i < arr.Len(); i++ {
		pie := arr.Index(i).Interface()

		switch vv := pie.(type) {
		case *stringStruct:
			buf.WriteString(vv.S)
		case string:
			buf.WriteString(vv)
		default:
			Panicerr("Unsupport type in Join")
		}
		if l--; l > 0 {
			buf.WriteString(s.S)
		}
	}
	s.S = buf.String()
	return s
}

func (s *stringStruct) Strip(characterMask ...string) *stringStruct {
	if len(characterMask) == 0 {
		s.S = strings.TrimSpace(s.S)
	} else {
		s.S = strings.Trim(s.S, characterMask[0])
	}
	return s
}

func (s *stringStruct) Split(sep ...string) []*stringStruct {
	var a []*stringStruct
	if len(sep) != 0 {
		for _, v := range strings.Split(s.S, sep[0]) {
			a = append(a, String(v).Strip())
		}
	} else {
		for _, v := range strings.Split(s.S, " ") {
			if String(v).Strip().Get() != "" {
				a = append(a, String(v).Strip())
			}
		}
	}

	return a
}

func (s *stringStruct) Utf8Split() []*stringStruct {
	var a []*stringStruct

	for _, v := range s.S {
		if String(string(v)).Strip().Get() != "" {
			a = append(a, String(string(v)).Strip())
		}
	}

	return a
}

func (s *stringStruct) Count(substr string) int {
	return strings.Count(s.S, substr)
}

func (s *stringStruct) EndsWith(substr string) (res bool) {
	if len(substr) <= len(s.S) && s.S[len(s.S)-len(substr):] == substr {
		res = true
	} else {
		res = false
	}
	return
}

func (s *stringStruct) StartsWith(substr string) (res bool) {
	if len(substr) <= len(s.S) && s.S[:len(substr)] == substr {
		res = true
	} else {
		res = false
	}
	return
}

func (s *stringStruct) Splitlines(strip ...bool) []*stringStruct {
	var a []*stringStruct
	tostrip := false
	if len(strip) != 0 {
		tostrip = strip[0]
	}
	for _, v := range strings.Split(s.S, "\n") {
		if tostrip {
			a = append(a, String(v).Strip())
		} else {
			a = append(a, String(v).RStrip("\r"))
		}
	}

	return a
}

func (s *stringStruct) In(str string) bool {
	return String(str).Index(s.S) != -1
}

func (s *stringStruct) LStrip(characterMask ...string) *stringStruct {
	if len(characterMask) == 0 {
		s.S = strings.TrimLeftFunc(s.S, unicode.IsSpace)
	}
	s.S = strings.TrimLeft(s.S, characterMask[0])
	return s
}

func (s *stringStruct) RStrip(characterMask ...string) *stringStruct {
	if len(characterMask) == 0 {
		s.S = strings.TrimRightFunc(s.S, unicode.IsSpace)
	}
	s.S = strings.TrimRight(s.S, characterMask[0])
	return s
}

func (ss *stringStruct) Isdigit() bool {
	str := ss.S
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

func (s *stringStruct) HasChinese() bool {
	var count int
	for _, v := range s.S {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

func (s *stringStruct) Filter(charts ...string) *stringStruct {
	var res string
	strb := []byte(s.S)
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
	s.S = res
	return s
}

func (s *stringStruct) RemoveHtmlTag() *stringStruct {
	s.S = striphtmltags.StripTags(s.S)
	return s
}

func (s *stringStruct) RemoveNonUTF8Character() *stringStruct {
	if !utf8.ValidString(s.S) {
		v := make([]rune, 0, len(s.S))
		for i, r := range s.S {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s.S[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		s.S = string(v)
	}
	return s
}

func (s *stringStruct) DetectLanguage() string {
	return whatlanggo.Detect(s.S).Lang.String()
}

func (s *stringStruct) IsAscii() bool {
	return utf8string.NewString(s.S).IsASCII()
}

func (s *stringStruct) RegexFindAll(pattern string, multiline ...bool) (res [][]*stringStruct) {
	if len(multiline) > 0 && multiline[0] {
		pattern = "(?s)" + pattern
	}
	r, err := regexp.Compile(pattern)
	Panicerr(err)

	for _, i := range r.FindAllStringSubmatch(s.S, -1) {
		var arr []*stringStruct
		for _, j := range i {
			arr = append(arr, String(j))
		}
		res = append(res, arr)
	}
	return
}

func (s *stringStruct) RegexReplace(pattern string, newstring string) *stringStruct {
	r, err := regexp.Compile(pattern)
	Panicerr(err)
	s.S = r.ReplaceAllString(s.S, newstring)
	return s
}
