package golanglibs

import (
	"math"
	"reflect"
	"strings"
)

type mathStruct struct {
	Abs          func(number float64) float64
	Sum          func(array interface{}) (sumresult float64)
	Average      func(array interface{}) (avgresult float64)
	DecimalToAny func(num, n int64) string
	AnyToDecimal func(num string, n int64) int64
}

var Math mathStruct

func init() {
	Math = mathStruct{
		Abs:          abs,
		Sum:          mathsum,
		Average:      mathaverage,
		DecimalToAny: decimalToAny,
		AnyToDecimal: anyToDecimal,
	}
}

var tenToAny map[int64]string = map[int64]string{
	0:  "0",
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "a",
	11: "b",
	12: "c",
	13: "d",
	14: "e",
	15: "f",
	16: "g",
	17: "h",
	18: "i",
	19: "j",
	20: "k",
	21: "l",
	22: "m",
	23: "n",
	24: "o",
	25: "p",
	26: "q",
	27: "r",
	28: "s",
	29: "t",
	30: "u",
	31: "v",
	32: "w",
	33: "x",
	34: "y",
	35: "z",
	36: "A",
	37: "B",
	38: "C",
	39: "D",
	40: "E",
	41: "F",
	42: "G",
	43: "H",
	44: "I",
	45: "J",
	46: "K",
	47: "L",
	48: "M",
	49: "N",
	50: "O",
	51: "P",
	52: "Q",
	53: "R",
	54: "S",
	55: "T",
	56: "U",
	57: "V",
	58: "W",
	59: "X",
	60: "Y",
	61: "Z"}

// 10进制转任意进制
func decimalToAny(num, n int64) string {
	newNumStr := ""
	var remainder int64
	var remainderString string
	for num != 0 {
		remainder = num % n
		if 76 > remainder && remainder > 9 {
			remainderString = tenToAny[remainder]
		} else {
			remainderString = Str(remainder)
		}
		newNumStr = remainderString + newNumStr
		num = num / n
	}
	return newNumStr
}

// 任意进制转10进制
func anyToDecimal(num string, n int64) int64 {
	var newNum float64
	newNum = 0.0
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := float64(func(in string) int64 {
			var result int64 = -1
			for k, v := range tenToAny {
				if in == v {
					result = k
				}
			}
			return result
		}(value))
		if tmp != -1 {
			newNum = newNum + tmp*math.Pow(float64(n), float64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return Int64(newNum)
}

func mathaverage(array interface{}) (avgresult float64) {
	arr := reflect.ValueOf(array)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		Panicerr("Invalid data type of param \"array\": Not an Array")
	}

	for i := 0; i < arr.Len(); i++ {
		avgresult += Float64(arr.Index(i).Interface())
	}

	avgresult = avgresult / Float64(arr.Len())

	return
}

func mathsum(array interface{}) (sumresult float64) {
	arr := reflect.ValueOf(array)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		Panicerr("Invalid data type of param \"array\": Not an Array")
	}

	for i := 0; i < arr.Len(); i++ {
		sumresult += Float64(arr.Index(i).Interface())
	}

	return
}

func abs(number float64) float64 {
	return math.Abs(number)
}
