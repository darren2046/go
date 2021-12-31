package golanglibs

import "regexp"

type reStruct struct {
	FindAll func(pattern string, text string, multiline ...bool) [][]*stringStruct
	Replace func(pattern string, newstring string, text string) string
}

var Re reStruct

func init() {
	Re = reStruct{
		FindAll: refindAll,
		Replace: reReplace,
	}
}

func refindAll(pattern string, text string, multiline ...bool) (res [][]*stringStruct) {
	if len(multiline) > 0 && multiline[0] {
		pattern = "(?s)" + pattern
	}
	r, err := regexp.Compile(pattern)
	Panicerr(err)

	for _, i := range r.FindAllStringSubmatch(text, -1) {
		var arr []*stringStruct
		for _, j := range i {
			arr = append(arr, String(j))
		}
		res = append(res, arr)
	}
	return
}

func reReplace(pattern string, newstring string, text string) string {
	r, err := regexp.Compile(pattern)
	Panicerr(err)
	return r.ReplaceAllString(text, newstring)
}
