package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

var bytesAes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func BcryptHash(value string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), cost)
	return string(bytes), err
}

func BcryptCheckHash(value, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}

func AesEncrypt(text string, aesSecret string) (string, error) {
	block, err := aes.NewCipher([]byte(aesSecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytesAes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return EncodeBase64(cipherText), nil
}

func AesDecrypt(text string, aesSecret string) (string, error) {
	block, err := aes.NewCipher([]byte(aesSecret))
	if err != nil {
		return "", err
	}
	cipherText, err := DecodeBase64(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, bytesAes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func EncodeBase64(b []byte) string {
	result := base64.StdEncoding.EncodeToString(b)
	return result
}
func DecodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}
