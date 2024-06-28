package utils

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	consts "github.com/lgrisa/lib/utils/const"
	"hash"
	"os"
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
		return consts.NULL, fmt.Errorf("[%v]: %+v", "rsa.SignPKCS1v15", err)
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

func Md5String(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

func ReadFileMd5(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return Md5String(data), nil
}

func DecodePublicKey(pemContent []byte) (publicKey *rsa.PublicKey, err error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("pem.Decode(%s)：pemContent decode error", pemContent)
	}
	switch block.Type {
	case "CERTIFICATE":
		pubKeyCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParseCertificate(%s)：%w", pemContent, err)
		}
		pubKey, ok := pubKeyCert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("公钥证书提取公钥出错 [%s]", pemContent)
		}
		publicKey = pubKey
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParsePKIXPublicKey(%s),err:%w", pemContent, err)
		}
		pubKey, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("公钥解析出错 [%s]", pemContent)
		}
		publicKey = pubKey
	case "RSA PUBLIC KEY":
		pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.ParsePKCS1PublicKey(%s)：%w", pemContent, err)
		}
		publicKey = pubKey
	}
	return publicKey, nil
}

func DecodePrivateKey(pemContent []byte) (privateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("pem.Decode(%s)：pemContent decode error", pemContent)
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		pk8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("私钥解析出错 [%s]", pemContent)
		}
		var ok bool
		privateKey, ok = pk8.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("私钥解析出错 [%s]", pemContent)
		}
	}
	return privateKey, nil
}
