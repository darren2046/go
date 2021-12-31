package golanglibs

import (
	"fmt"
	"reflect"
	"unicode/utf8"

	reprlib "github.com/alecthomas/repr"
	"github.com/k0kubun/pp"
)

// Convert ASCII number to char
func Chr(ascii int) string {
	return string(rune(ascii))
}

//  Convert char to ASCII number
func Ord(char string) int {
	r, _ := utf8.DecodeRune([]byte(char))
	return int(r)
}

func Repr(obj interface{}) string {
	return reprlib.String(obj)
}

func Print(data ...interface{}) int {
	scheme := pp.ColorScheme{
		Bool:            pp.Cyan | pp.Bold,
		Integer:         pp.Blue | pp.Bold,
		Float:           pp.Magenta | pp.Bold,
		String:          pp.Green,
		StringQuotation: pp.Green | pp.Bold,
		EscapedChar:     pp.Magenta | pp.Bold,
		FieldName:       pp.Yellow,
		PointerAdress:   pp.Blue | pp.Bold,
		Nil:             pp.Cyan | pp.Bold,
		Time:            pp.Blue | pp.Bold,
		StructName:      pp.Green | pp.Bold,
		ObjectLength:    pp.Blue,
	}

	pp.SetColorScheme(scheme)

	num, err := pp.Println(data...)
	Panicerr(err)
	return num
}

func Printf(format string, data ...interface{}) int {
	scheme := pp.ColorScheme{
		Bool:            pp.Cyan | pp.Bold,
		Integer:         pp.Blue | pp.Bold,
		Float:           pp.Magenta | pp.Bold,
		String:          pp.Green,
		StringQuotation: pp.Green | pp.Bold,
		EscapedChar:     pp.Magenta | pp.Bold,
		FieldName:       pp.Yellow,
		PointerAdress:   pp.Blue | pp.Bold,
		Nil:             pp.Cyan | pp.Bold,
		Time:            pp.Blue | pp.Bold,
		StructName:      pp.Green | pp.Bold,
		ObjectLength:    pp.Blue,
	}

	pp.SetColorScheme(scheme)

	num, err := pp.Printf(format, data...)
	Panicerr(err)
	return num
}

func Sprint(data ...interface{}) string {
	scheme := pp.ColorScheme{
		Bool:            pp.Cyan | pp.Bold,
		Integer:         pp.Blue | pp.Bold,
		Float:           pp.Magenta | pp.Bold,
		String:          pp.Green,
		StringQuotation: pp.Green | pp.Bold,
		EscapedChar:     pp.Magenta | pp.Bold,
		FieldName:       pp.Yellow,
		PointerAdress:   pp.Blue | pp.Bold,
		Nil:             pp.Cyan | pp.Bold,
		Time:            pp.Blue | pp.Bold,
		StructName:      pp.Green | pp.Bold,
		ObjectLength:    pp.Blue,
	}

	//pp.ColoringEnabled = false
	pp.SetColorScheme(scheme)

	return String(pp.Sprintln(data...)).Strip().Get()
}

func Range(num ...int) []int {
	if len(num) != 1 && len(num) != 2 {
		Panicerr("需要1～2个参数")
	}
	var a []int
	var start int
	var end int
	if len(num) == 1 {
		start = 0
		end = num[0]
	} else {
		start = num[0]
		end = num[1]
	}
	for i := start; i < end; i++ {
		a = append(a, i)
	}
	return a
}

func Typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func Input(prompt string, defaultValue ...string) *stringStruct {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		if len(defaultValue) != 0 {
			input = defaultValue[0]
		}
	}
	return String(input)
}
