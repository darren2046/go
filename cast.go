package golanglibs

import "github.com/spf13/cast"

func str(v interface{}) string {
	return cast.ToString(v)
}

func toFloat64(v interface{}) float64 {
	return cast.ToFloat64(v)
}

func toInt64(v interface{}) int64 {
	return cast.ToInt64(v)
}
