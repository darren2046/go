package golanglibs

import "golang.org/x/net/html"

type htmlStruct struct {
	Encode func(str string) string
	Decode func(str string) string
}

var Html htmlStruct

func init() {
	Html = htmlStruct{
		Encode: htmlencode,
		Decode: htmldecode,
	}
}

//  将字符转换为 HTML 转义字符
func htmlencode(str string) string {
	return html.EscapeString(str)
}

//  Convert HTML entities to their corresponding characters
func htmldecode(str string) string {
	return html.UnescapeString(str)
}
