package golanglibs

import "encoding/base64"

type base64Struct struct {
	Encode func(str string) string
	Decode func(str string) string
}

var Base64 base64Struct

func init() {
	Base64 = base64Struct{
		Encode: base64Encode,
		Decode: base64Decode,
	}
}

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) string {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	Panicerr(err)
	return string(data)
}
