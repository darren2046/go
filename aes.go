package golanglibs

import (
	"crypto/aes"
	"crypto/cipher"
)

type aesStruct struct {
	key []byte
}

// Key为string
func getAES(key interface{}) *aesStruct {
	return &aesStruct{key: []byte(Str(key))}
}

// plaintext is string
func (a aesStruct) Encrypt(plaintext interface{}) *stringStruct {
	block, err := aes.NewCipher(a.key)
	Panicerr(err)
	ciphertext := make([]byte, aes.BlockSize+len(Str(plaintext)))
	iv := ciphertext[:aes.BlockSize]
	Panicerr(err)
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(Str(plaintext)))
	return String(string(ciphertext))
	//return hex.EncodeToString(ciphertext)
}

// encryptedText is string
func (a aesStruct) Decrypt(encryptedText interface{}) *stringStruct {
	ciphertext := []byte(Str(encryptedText))
	//ciphertext, err := hex.DecodeString(d)
	block, err := aes.NewCipher(a.key)
	Panicerr(err)
	if len(ciphertext) < aes.BlockSize {
		panic(newerr("ciphertext too short"))
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return String(string(ciphertext))
}
