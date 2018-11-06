package sessionHandler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

const key = "rRTNx^29L2.^?6mC94Cu%W?8wBwKaHoz]zFT7[wB"

func hashTo32Bytes(input string) []byte {
	data := sha256.Sum256([]byte(input))
	return data[0:]
}

func getKey() string {
	return key
}

func encryptString(plainText string, keyString string) (cipherTextString string, err error) {
	key := hashTo32Bytes(keyString)
	encrypted, err := encryptAES(key, []byte(plainText))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func encryptAES(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(encrypted, data)
	return output, nil
}

func decryptString(cryptoText string, keyString string) (plainTextString string, err error) {
	encrypted, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(encrypted) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short. It decodes to %v bytes but the minimum length is 16", len(encrypted))
	}

	decrypted, err := decryptAES(hashTo32Bytes(keyString), encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func decryptAES(key, data []byte) ([]byte, error) {
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}
