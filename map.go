package golanglibs

import "reflect"

type mapStruct struct {
	mmap reflect.Value
}

func Map(mmap interface{}) *mapStruct {
	mmmap := reflect.ValueOf(mmap)
	if mmmap.Kind() != reflect.Map {
		Panicerr("Invalid data type of param \"Map\": " + Typeof(mmap))
	}
	return &mapStruct{mmap: mmmap}
}

func (m *mapStruct) Has(Key interface{}) bool {
	// Print(Typeof(Key))
	for _, v := range m.mmap.MapKeys() {
		if v.Interface() == Key {
			return true
		}
	}
	return false
}

func (m *mapStruct) Keys() (keys []string) {
	for _, v := range m.mmap.MapKeys() {
		keys = append(keys, Str(v.Interface()))
	}
	return
}
