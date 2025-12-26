package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func Aes256Encode(plainText string, iv string, aesKey string) (string, error) {
	byteAesKey, err := hex.DecodeString(aesKey)
	if err != nil {
		return "", err
	}

	bIV, err := hex.DecodeString(iv)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(byteAesKey)
	if err != nil {
		return "", err
	}

	plaintextBytes := PKCS7Padding([]byte(plainText), aes.BlockSize)
	ciphertext := make([]byte, len(plaintextBytes))

	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, plaintextBytes)

	return hex.EncodeToString(ciphertext), nil
}

func Aes256Decode(cipherText string, iv string, aesKey string) (string, error) {
	byteAesKey, err := hex.DecodeString(aesKey)
	if err != nil {
		return "", err
	}

	bIV, err := hex.DecodeString(iv)
	if err != nil {
		return "", err
	}

	cipherTextDecoded, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(byteAesKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))
	return string(bytes.TrimSpace(cipherTextDecoded)), nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
