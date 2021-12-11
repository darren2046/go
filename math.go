package golanglibs

import (
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
)

type mathStruct struct {
	Int2ip       func(ipnr int64) string
	Ip2int       func(ipnr string) int64
	Abs          func(number float64) float64
	Sum          func(array interface{}) (sumresult float64)
	Average      func(array interface{}) (avgresult float64)
	DecimalToAny func(num, n int64) string
	AnyToDecimal func(num string, n int64) int64
	IpInNet      func(ip string, Net string, mask ...string) bool
}

var Math mathStruct

func init() {
	Math = mathStruct{
		Int2ip:       int2ip,
		Ip2int:       ip2int,
		Abs:          abs,
		Sum:          mathsum,
		Average:      mathaverage,
		DecimalToAny: decimalToAny,
		AnyToDecimal: anyToDecimal,
		IpInNet:      ipInNet,
	}
}

func ipInNet(ip string, Net string, mask ...string) bool {
	if len(mask) != 0 {
		ip := net.ParseIP(mask[0])
		addr := ip.To4()
		cidrsuffix, _ := net.IPv4Mask(addr[0], addr[1], addr[2], addr[3]).Size()
		Net = Net + "/" + Str(cidrsuffix)
	}

	_, ipnetA, _ := net.ParseCIDR(Net)
	ipB := net.ParseIP(ip)

	if ipnetA.Contains(ipB) {
		return true
	} else {
		return false
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
	new_num_str := ""
	var remainder int64
	var remainder_string string
	for num != 0 {
		remainder = num % n
		if 76 > remainder && remainder > 9 {
			remainder_string = tenToAny[remainder]
		} else {
			remainder_string = Str(remainder)
		}
		new_num_str = remainder_string + new_num_str
		num = num / n
	}
	return new_num_str
}

// 任意进制转10进制
func anyToDecimal(num string, n int64) int64 {
	var new_num float64
	new_num = 0.0
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
			new_num = new_num + tmp*math.Pow(float64(n), float64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return Int64(new_num)
}

func mathaverage(array interface{}) (avgresult float64) {
	arr := reflect.ValueOf(array)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		panicerr("Invalid data type of param \"array\": Not an Array")
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
		panicerr("Invalid data type of param \"array\": Not an Array")
	}

	for i := 0; i < arr.Len(); i++ {
		sumresult += Float64(arr.Index(i).Interface())
	}

	return
}

func abs(number float64) float64 {
	return math.Abs(number)
}

func int2ip(ipnr int64) string {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

func ip2int(ipnr string) int64 {
	bits := strings.Split(ipnr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}
