package golanglibs

import (
	"crypto/cipher"
	"math/rand"

	"golang.org/x/crypto/chacha20poly1305"
)

type chacha20poly1305Struct struct {
	aead cipher.AEAD
}

// key length need to be 32 bytes
func getChacha20poly1305(key string) *chacha20poly1305Struct {
	if len(key) != 32 {
		Panicerr("key length need to be 32 bytes")
	}
	aead, err := chacha20poly1305.NewX([]byte(key))
	Panicerr(err)

	return &chacha20poly1305Struct{
		aead: aead,
	}
}

func (m *chacha20poly1305Struct) Encrypt(plantext string) (ciphertext string) {
	msg := []byte(plantext)
	nonce := make([]byte, m.aead.NonceSize(), m.aead.NonceSize()+len(msg)+m.aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg := m.aead.Seal(nonce, nonce, msg, nil)
	ciphertext = string(encryptedMsg)
	return
}

func (m *chacha20poly1305Struct) Decrypt(ciphertext string) (plaintext string) {
	encryptedMsg := []byte(ciphertext)

	if len(encryptedMsg) < m.aead.NonceSize() {
		panic("ciphertext too short")
	}

	nonce, ciphertextbyte := encryptedMsg[:m.aead.NonceSize()], encryptedMsg[m.aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintextbyte, err := m.aead.Open(nil, nonce, ciphertextbyte, nil)
	if err != nil {
		panic(err)
	}
	plaintext = string(plaintextbyte)

	return
}
