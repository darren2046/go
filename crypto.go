package golanglibs

type cryptoStruct struct {
	Xor func(data, key string) string
	Aes func(key string) *aesStruct
}

var Crypto cryptoStruct

func init() {
	Crypto = cryptoStruct{
		Xor: xor,
		Aes: getAES,
	}
}

func xor(data, key string) string {
	var output []byte
	keyarr := []byte(key)
	dataarr := []byte(data)
	kL := len(keyarr)
	for i := range dataarr {
		output = append(output, data[i]^key[i%kL])
	}
	return string(output)
}
