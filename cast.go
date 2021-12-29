package golanglibs

import "github.com/spf13/cast"

func Str(v interface{}) string {
	return cast.ToString(v)
}

func Float64(v interface{}) float64 {
	return cast.ToFloat64(v)
}

func Int64(v interface{}) int64 {
	return cast.ToInt64(v)
}

func Int(v interface{}) int {
	return cast.ToInt(v)
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
