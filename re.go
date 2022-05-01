package golanglibs

import "regexp"

type reStruct struct {
	FindAll func(pattern string, text string, multiline ...bool) [][]*StringStruct
	Replace func(pattern string, newstring string, text string) string
}

type reCompiledPatternStruct struct {
	re   *regexp.Regexp
	time float64
}

var reCompiledPatternMap map[string]*reCompiledPatternStruct

var Re reStruct

var reRWLock *RWLockStruct

func init() {
	reRWLock = getRWLock()

	reCompiledPatternMap = make(map[string]*reCompiledPatternStruct)

	Re = reStruct{
		FindAll: refindAll,
		Replace: reReplace,
	}

	go func() {
		for {
			sleep(120)
			reRWLock.WAcquire()
			// for 循环的时候这个map不能修改否则就是
			for k, v := range reCompiledPatternMap {
				// 120秒清一次正则的缓存
				if v.time < timeNowInTimestamp()-120 {
					delete(reCompiledPatternMap, k)
				}
			}
			reRWLock.WRelease()
		}
	}()
}

func refindAll(pattern string, text string, multiline ...bool) (res [][]*StringStruct) {
	if len(multiline) > 0 && multiline[0] {
		pattern = "(?s)" + pattern
	}

	reRWLock.RAcquire()
	defer reRWLock.RRelease()

	if _, ok := reCompiledPatternMap[pattern]; !ok {
		r, err := regexp.Compile(pattern)
		Panicerr(err)
		reCompiledPatternMap[pattern] = &reCompiledPatternStruct{
			re:   r,
			time: Time.Now(),
		}
	}

	for _, i := range reCompiledPatternMap[pattern].re.FindAllStringSubmatch(text, -1) {
		var arr []*StringStruct
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
