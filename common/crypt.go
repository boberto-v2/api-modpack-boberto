package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

var bytesAes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Encrypt(text string, aesSecret string) (string, error) {
	block, err := aes.NewCipher([]byte(aesSecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytesAes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text string, aesSecret string) (string, error) {
	block, err := aes.NewCipher([]byte(aesSecret))
	if err != nil {
		return "", err
	}
	cipherText, err := decode(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, bytesAes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}
