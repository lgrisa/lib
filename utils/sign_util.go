package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
)

const (
	RSA  = "RSA"
	RSA2 = "RSA2"
)

func GetRsaSign(body, signType string, privateKey *rsa.PrivateKey) (sign string, err error) {
	var (
		h              hash.Hash
		hashes         crypto.Hash
		encryptedBytes []byte
	)

	switch signType {
	case RSA:
		h = sha1.New()
		hashes = crypto.SHA1
	case RSA2:
		h = sha256.New()
		hashes = crypto.SHA256
	default:
		h = sha256.New()
		hashes = crypto.SHA256
	}

	if _, err = h.Write([]byte(body)); err != nil {
		return
	}
	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, privateKey, hashes, h.Sum(nil)); err != nil {
		return NULL, fmt.Errorf("[%v]: %+v", "rsa.SignPKCS1v15", err)
	}

	sign = base64.StdEncoding.EncodeToString(encryptedBytes)

	return
}

func VerifySign(signData, sign, signType, alipayPublicKey string) (err error) {
	var (
		h     hash.Hash
		hashs crypto.Hash
	)
	publicKey, err := DecodePublicKey([]byte(alipayPublicKey))
	if err != nil {
		return err
	}
	signBytes, _ := base64.StdEncoding.DecodeString(sign)

	switch signType {
	case RSA:
		hashs = crypto.SHA1
	case RSA2:
		hashs = crypto.SHA256
	default:
		hashs = crypto.SHA256
	}
	h = hashs.New()
	h.Write([]byte(signData))
	if err = rsa.VerifyPKCS1v15(publicKey, hashs, h.Sum(nil), signBytes); err != nil {
		return fmt.Errorf("[%v]: %v", "VerifySign", err)
	}
	return nil
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
