package golanglibs

import "reflect"

type mapStruct struct {
	mmap interface{}
}

func Map(mmap interface{}) *mapStruct {
	return &mapStruct{mmap: mmap}
}

func (m *mapStruct) Has(Key interface{}) bool {
	arr := reflect.ValueOf(Map)
	if arr.Kind() != reflect.Map {
		panicerr("Invalid data type of param \"Map\": Not a Map")
	}
	for _, v := range arr.MapKeys() {
		if v.Interface() == Key {
			return true
		}
	}
	return false
}

func (m *mapStruct) Keys() (keys []string) {
	arr := reflect.ValueOf(m.mmap)
	if arr.Kind() != reflect.Map {
		panicerr("Invalid data type of param \"Map\": Not a Map")
	}
	for _, v := range arr.MapKeys() {
		keys = append(keys, Str(v.Interface()))
	}
	return
}
