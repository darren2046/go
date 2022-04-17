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

	"github.com/denisbrodbeck/striphtmltags"
	"github.com/pemistahl/lingua-go"
	"golang.org/x/exp/utf8string"
)

var languageDetector lingua.LanguageDetector

func init() {
	languages := []lingua.Language{
		lingua.English,
		lingua.French,
		lingua.German,
		lingua.Spanish,
		lingua.Chinese,
		lingua.Hindi,
		lingua.Italian,
		lingua.Korean,
		lingua.Russian,
		lingua.Indonesian,
		lingua.Arabic,
		lingua.Turkish,
		lingua.Polish,
		lingua.Swedish,
		lingua.Thai,
		lingua.Malay,
		lingua.Vietnamese,
	}

	languageDetector = lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()
}

type StringStruct struct {
	S string
}

func String(s interface{}) *StringStruct {
	return &StringStruct{S: Str(s)}
}

// Return the final string
func (s *StringStruct) Get() string {
	return s.S
}

// Return the substring of the string
func (s *StringStruct) Sub(start, end int) *StringStruct {
	s.S = s.sub(start, end)
	return s
}

func (s *StringStruct) Has(substr string) bool {
	return strings.Contains(s.S, substr)
}

func (s *StringStruct) sub(start, end int) string {
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
func (s *StringStruct) Len() int {
	return len(s.S)
}

// Reverse the string
func (s *StringStruct) Reverse() *StringStruct {
	runes := []rune(s.S)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	s.S = string(runes)
	return s
}

// Split the string with specific length
func (s *StringStruct) Chunk(length int) (res []*StringStruct) {
	start := 0
	for s.Len() > start {
		substr := s.sub(start, start+length)
		res = append(res, String(substr))
		start = start + length
	}
	return
}

// Calculate the length of a UTF-8 string
func (s *StringStruct) Utf8Len() int {
	return utf8.RuneCountInString(s.S)
}

// Repeat a string
func (s *StringStruct) Repeat(count int) *StringStruct {
	s.S = strings.Repeat(s.S, count)
	return s
}

// Shuffle a string
func (s *StringStruct) Shuffle() *StringStruct {
	runes := []rune(s.S)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ss := make([]rune, len(runes))
	for i, v := range r.Perm(len(runes)) {
		ss[i] = runes[v]
	}
	s.S = string(ss)
	return s
}

func (s *StringStruct) Index(substr string) int {
	return strings.Index(s.S, substr)
}

func (s *StringStruct) Replace(oldText, newText interface{}) *StringStruct {
	s.S = strings.ReplaceAll(s.S, Str(oldText), Str(newText))
	return s
}

func (s *StringStruct) Upper() *StringStruct {
	s.S = strings.ToUpper(s.S)
	return s
}

func (s *StringStruct) Lower() *StringStruct {
	s.S = strings.ToLower(s.S)
	return s
}

// pieces can be []string or []*StringStruct
func (s *StringStruct) Join(pieces interface{}) *StringStruct {
	var buf bytes.Buffer
	arr := reflect.ValueOf(pieces)
	l := arr.Len()

	for i := 0; i < arr.Len(); i++ {
		pie := arr.Index(i).Interface()

		switch vv := pie.(type) {
		case *StringStruct:
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

func (s *StringStruct) Strip(characterMask ...string) *StringStruct {
	if len(characterMask) == 0 {
		s.S = strings.TrimSpace(s.S)
	} else {
		s.S = strings.Trim(s.S, characterMask[0])
	}
	return s
}

func (s *StringStruct) Split(sep ...string) []*StringStruct {
	var a []*StringStruct
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

func (s *StringStruct) Utf8Split() []*StringStruct {
	var a []*StringStruct

	for _, v := range s.S {
		if String(string(v)).Strip().Get() != "" {
			a = append(a, String(string(v)).Strip())
		}
	}

	return a
}

func (s *StringStruct) Count(substr string) int {
	return strings.Count(s.S, substr)
}

func (s *StringStruct) EndsWith(substr string) (res bool) {
	if len(substr) <= len(s.S) && s.S[len(s.S)-len(substr):] == substr {
		res = true
	} else {
		res = false
	}
	return
}

func (s *StringStruct) StartsWith(substr string) (res bool) {
	if len(substr) <= len(s.S) && s.S[:len(substr)] == substr {
		res = true
	} else {
		res = false
	}
	return
}

func (s *StringStruct) Splitlines(strip ...bool) []*StringStruct {
	var a []*StringStruct
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

func (s *StringStruct) In(str interface{}) bool {
	switch vv := str.(type) {
	case *StringStruct:
		return vv.Index(s.S) != -1
	case string:
		return String(vv).Index(s.S) != -1
	default:
		Panicerr("Unsupport type in In()")
	}
	return false
}

func (s *StringStruct) LStrip(characterMask ...string) *StringStruct {
	if len(characterMask) == 0 {
		s.S = strings.TrimLeftFunc(s.S, unicode.IsSpace)
	}
	s.S = strings.TrimLeft(s.S, characterMask[0])
	return s
}

func (s *StringStruct) RStrip(characterMask ...string) *StringStruct {
	if len(characterMask) == 0 {
		s.S = strings.TrimRightFunc(s.S, unicode.IsSpace)
	}
	s.S = strings.TrimRight(s.S, characterMask[0])
	return s
}

func (ss *StringStruct) Isdigit() bool {
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

func (s *StringStruct) HasChinese() bool {
	var count int
	for _, v := range s.S {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

func (s *StringStruct) Filter(charts ...string) *StringStruct {
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

func (s *StringStruct) RemoveHtmlTag() *StringStruct {
	s.S = striphtmltags.StripTags(s.S)
	return s
}

func (s *StringStruct) RemoveNonUTF8Character() *StringStruct {
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

// 返回字符串的语言, 为了资源使用考虑，只支持2021年GDP前30的国家。
func (s *StringStruct) Language() string {
	if language, exists := languageDetector.DetectLanguageOf(s.S); exists {
		return Str(language)
	} else {
		return "Unknown"
	}
}

func (s *StringStruct) IsAscii() bool {
	return utf8string.NewString(s.S).IsASCII()
}

func (s *StringStruct) RegexFindAll(pattern string, multiline ...bool) (res [][]*StringStruct) {
	if len(multiline) > 0 && multiline[0] {
		pattern = "(?s)" + pattern
	}
	r, err := regexp.Compile(pattern)
	Panicerr(err)

	for _, i := range r.FindAllStringSubmatch(s.S, -1) {
		var arr []*StringStruct
		for _, j := range i {
			arr = append(arr, String(j))
		}
		res = append(res, arr)
	}
	return
}

func (s *StringStruct) RegexReplace(pattern string, newstring string) *StringStruct {
	r, err := regexp.Compile(pattern)
	Panicerr(err)
	s.S = r.ReplaceAllString(s.S, newstring)
	return s
}

func (s *StringStruct) JsonXPath() *xpathJsonStruct {
	return getXPathJson(s.S)
}

func (s *StringStruct) XPath() *xpathStruct {
	return getXPath(s.S)
}

// return the len of longest string both in str1 and str2 and the positions in str1 and str2
func SimilarStr(str1 []rune, str2 []rune) (int, int, int) {
	var maxLen, tmp, pos1, pos2 = 0, 0, 0, 0
	len1, len2 := len(str1), len(str2)

	for p := 0; p < len1; p++ {
		for q := 0; q < len2; q++ {
			tmp = 0
			for p+tmp < len1 && q+tmp < len2 && str1[p+tmp] == str2[q+tmp] {
				tmp++
			}
			if tmp > maxLen {
				maxLen, pos1, pos2 = tmp, p, q
			}
		}

	}

	return maxLen, pos1, pos2
}

// return the total length of longest string both in str1 and str2
func SimilarChar(str1 []rune, str2 []rune) int {
	maxLen, pos1, pos2 := SimilarStr(str1, str2)
	total := maxLen

	if maxLen != 0 {
		if pos1 > 0 && pos2 > 0 {
			total += SimilarChar(str1[:pos1], str2[:pos2])
		}
		if pos1+maxLen < len(str1) && pos2+maxLen < len(str2) {
			total += SimilarChar(str1[pos1+maxLen:], str2[pos2+maxLen:])
		}
	}

	return total
}

// return a int value in [0, 100], which stands for match level
func (s *StringStruct) Similar(text string) int {
	txt1, txt2 := []rune(s.S), []rune(text)
	if len(txt1) == 0 || len(txt2) == 0 {
		return 0
	}
	return SimilarChar(txt1, txt2) * 200 / (len(txt1) + len(txt2))
}
