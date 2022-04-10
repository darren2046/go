package golanglibs

import (
	"reflect"

	"github.com/spf13/cast"
)

func Str(v interface{}) (res string) {
	switch vv := v.(type) {
	case *StringStruct:
		res = cast.ToString(vv.S)
	default:
		res = cast.ToString(v)
	}
	return
}

func Float64(v interface{}) float64 {
	return cast.ToFloat64(v)
}

func StringMap(v interface{}) map[string]interface{} {
	return cast.ToStringMap(v)
}

func Int64(v interface{}) int64 {
	return cast.ToInt64(v)
}

func Int(v interface{}) int {
	return cast.ToInt(v)
}

func Int32(v interface{}) int32 {
	return cast.ToInt32(v)
}

func Bool(v interface{}) bool {
	return cast.ToBool(v)
}

func Uint(v interface{}) uint {
	return cast.ToUint(v)
}

func Uint16(v interface{}) uint16 {
	return cast.ToUint16(v)
}

func Uint32(v interface{}) uint32 {
	return cast.ToUint32(v)
}

func InterfaceArray(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
