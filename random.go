package golanglibs

import (
	"math/rand"
	"reflect"
	"time"
)

type randomStruct struct {
	Int    func(min, max int64) int64
	Choice func(array interface{}) interface{}
	String func(length int, charset ...string) string
}

var Random randomStruct

func init() {
	Random = randomStruct{
		Int:    randint,
		Choice: randomChoice,
		String: randomStr,
	}
}

func randomStr(length int, charset ...string) string {
	var str string
	if len(charset) != 0 {
		str = charset[0]
	} else {
		str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = str[seededRand.Intn(len(str))]
	}
	return string(b)
}

func randomChoice(array interface{}) interface{} {
	rand.Seed(Int64(Time.Now() * 1000000))
	arr := reflect.ValueOf(array)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		Panicerr("Invalid data type of param \"array\": Not an Array")
	}

	return arr.Index(rand.Intn(arr.Len())).Interface()
}

func randint(min, max int64) int64 {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if max > 9223372036854775807 {
		panic("max: max can not be greater than 9223372036854775807")
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max+1-min) + min
}
