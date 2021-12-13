package golanglibs

import (
	"crypto/aes"
	"crypto/cipher"
)

type aesStruct struct {
	key []byte
}

func getAES(key string) *aesStruct {
	return &aesStruct{key: []byte(key)}
}

func (a aesStruct) Encrypt(plaintext string) string {
	block, err := aes.NewCipher(a.key)
	Panicerr(err)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	Panicerr(err)
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return string(ciphertext)
	//return hex.EncodeToString(ciphertext)
}

func (a aesStruct) Decrypt(d string) string {
	ciphertext := []byte(d)
	//ciphertext, err := hex.DecodeString(d)
	block, err := aes.NewCipher(a.key)
	Panicerr(err)
	if len(ciphertext) < aes.BlockSize {
		panic(newerr("ciphertext too short"))
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext)
}
