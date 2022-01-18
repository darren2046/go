package golanglibs

import (
	"testing"
)

func TestChacha20Poly1305(t *testing.T) {
	cha := Crypto.ChaCha20Poly1305("6UrvN36kG1ZEECGJgJEeYAyZfdtBfV00")
	s := "KeySize is the size of the key used by this AEAD, in bytes."
	if cha.Decrypt(cha.Encrypt(s)) != s {
		t.Error("ChaCha20Poly1305 error")
	}
}
