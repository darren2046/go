package golanglibs

import (
	"bytes"
	"encoding/binary"
)

type binaryStruct struct {
	Map2bin func(m map[string]string) string
	Bin2map func(s string) (res map[string]string)
}

var Binary binaryStruct

func init() {
	Binary = binaryStruct{
		Map2bin: map2bin,
		Bin2map: bin2map,
	}
}

func map2bin(m map[string]string) (res string) {
	btlen := make([]byte, 4)
	binary.LittleEndian.PutUint32(btlen, uint32(len(m)))
	res += string(btlen) // 元素个数

	for k, v := range m {
		binary.LittleEndian.PutUint32(btlen, uint32(len(k)))
		res += string(btlen) // key长度
		res += k             // key

		binary.LittleEndian.PutUint32(btlen, uint32(len(v)))
		res += string(btlen) // value长度
		res += v             // value
	}
	return
}

func bin2map(s string) (res map[string]string) {
	res = make(map[string]string)

	btlen := make([]byte, 4)

	ss := bytes.NewBuffer([]byte(s))

	ss.Read(btlen)
	tlen := int(binary.LittleEndian.Uint32(btlen)) // 元素个数

	for range Range(tlen) {
		ss.Read(btlen)
		tlen := int(binary.LittleEndian.Uint32(btlen)) // key长度
		key := make([]byte, tlen)
		ss.Read(key) // 读取key

		ss.Read(btlen)
		tlen = int(binary.LittleEndian.Uint32(btlen)) // value长度
		value := make([]byte, tlen)
		ss.Read(value) // 读取key

		res[string(key)] = string(value)
	}
	return
}
