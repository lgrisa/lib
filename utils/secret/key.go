package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var (
	ErrInvalidAesKey = errors.New("Invalid AES Key")
	ErrInvalidData   = errors.New("Invalid Data")
)

func Unencrypt(unencryptKey cipher.Block, encrypted []byte) []byte {
	if len(encrypted) <= 16 {
		return nil
	}

	iv := encrypted[:16]
	data := encrypted[16:]

	stream := cipher.NewCFBDecrypter(unencryptKey, iv)
	stream.XORKeyStream(data, data)
	return data
}

func UnencryptBytes(aesKey, encrypted []byte) ([]byte, error) {
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, ErrInvalidAesKey
	}
	return Unencrypt(c, encrypted), nil
}

func UnencryptString(aesKey, encrypted string) (string, error) {
	b, err := base64.RawStdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", ErrInvalidData
	}
	data, err := UnencryptBytes([]byte(aesKey), b)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
