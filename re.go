package golanglibs

import "regexp"

type reStruct struct {
	FindAll func(pattern string, text string, multiline ...bool) [][]string
	Replace func(pattern string, newstring string, text string) string
}

var Re reStruct

func init() {
	Re = reStruct{
		FindAll: refindAll,
		Replace: reReplace,
	}
}

func refindAll(pattern string, text string, multiline ...bool) [][]string {
	if len(multiline) > 0 && multiline[0] {
		pattern = "(?s)" + pattern
	}
	r, err := regexp.Compile(pattern)
	panicerr(err)
	return r.FindAllStringSubmatch(text, -1)
}

func reReplace(pattern string, newstring string, text string) string {
	r, err := regexp.Compile(pattern)
	panicerr(err)
	return r.ReplaceAllString(text, newstring)
}
