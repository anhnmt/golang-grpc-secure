package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// PKCS7Padding PKCS7Padding
func PKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding PKCS7UnPadding
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// EncryptAES EncryptAES
func EncryptAES(data, key string) (string, error) {
	plaintext := PKCS7Padding([]byte(data))

	block, err := aes.NewCipher([]byte(key[:32]))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(plaintext))

	// Create secret IV
	iv := [16]byte{}
	mode := cipher.NewCBCEncrypter(block, iv[:])
	mode.CryptBlocks(ciphertext, plaintext)

	encrypted := base64.StdEncoding.EncodeToString(ciphertext)

	return encrypted, nil
}

// DecryptAES DecryptAES
func DecryptAES(data, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key[:32]))
	if err != nil {
		return "", err
	}

	// Create secret IV
	iv := [16]byte{}
	mode := cipher.NewCBCDecrypter(block, iv[:])
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted = PKCS7UnPadding(decrypted)

	return string(decrypted), nil
}

// EncryptRSA encrypt RSA
func EncryptRSA(data string, publicKey []byte) (string, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the public key")
	}

	rsaPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(data))
	if err != nil {
		return "", err
	}

	encrypted := base64.StdEncoding.EncodeToString(ciphertext)
	return encrypted, nil
}

// DecryptRSA decrypt RSA
func DecryptRSA(data, privateKey []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the private key")
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, data)
	if err != nil {
		return "", err
	}

	fmt.Println("Plaintext:", string(plaintext))
	return string(plaintext), nil
}

// GenerateCheckSum GenerateCheckSum
func GenerateCheckSum(phone, key, msgType string, now int64) string {
	data := fmt.Sprintf("%s%s%s%v%s", phone, fmt.Sprintf("%d000000", now), msgType, float64(now)/1e12, "E12")

	checkSum, _ := EncryptAES(data, key)
	return checkSum
}
