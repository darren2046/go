package golanglibs

type cryptoStruct struct {
	Xor func(data, key interface{}) interface{} // data和key都是string
	Aes func(key interface{}) *aesStruct        // Key为string
}

var Crypto cryptoStruct

func init() {
	Crypto = cryptoStruct{
		Xor: xor,
		Aes: getAES,
	}
}

func xor(data, key interface{}) interface{} {
	var output []byte
	keystr := Str(key)
	datastr := Str(data)
	keyarr := []byte(keystr)
	dataarr := []byte(datastr)
	kL := len(keyarr)
	for i := range dataarr {
		output = append(output, datastr[i]^keystr[i%kL])
	}
	return string(output)
}
