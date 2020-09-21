package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// InArray to get if value is in array.
func InArray(arr []string, v string) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}

// Encrypt to encrypt data with key.
func Encrypt(str, key string) string {
	for len(key) < 32 {
		key += key
	}

	if len(key) > 32 {
		key = key[0:32]
	}

	block, _ := aes.NewCipher([]byte(key))
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	ciphertextByte := gcm.Seal(nonce, nonce, []byte(str), nil)

	return base64.StdEncoding.EncodeToString(ciphertextByte)
}

// Decrypt to decrypt encrypted data with key.
func Decrypt(str, key string) string {
	for len(key) < 32 {
		key += key
	}

	if len(key) > 32 {
		key = key[0:32]
	}

	block, _ := aes.NewCipher([]byte(key))
	gcm, _ := cipher.NewGCM(block)

	nonceSize := gcm.NonceSize()
	ciphertextByte, _ := base64.StdEncoding.DecodeString(str)
	nonce, ciphertextByteClean := ciphertextByte[:nonceSize], ciphertextByte[nonceSize:]
	plaintextByte, _ := gcm.Open(nil, nonce, ciphertextByteClean, nil)

	return string(plaintextByte)
}

// RandomStr to generate random string.
func RandomStr(l int) string {
	bytes := make([]byte, l)
	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes)
}
