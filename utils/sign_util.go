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
	"github.com/pkg/errors"
	"hash"
	"os"
	"strings"
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

const (
	PKCS1 PKCSType = 1 // 非java适用
	PKCS8 PKCSType = 2 // java适用
)

type PKCSType uint8

// FormatAlipayPrivateKey 格式化支付宝普通应用秘钥
func FormatAlipayPrivateKey(privateKey string) (pKey string) {
	var buffer strings.Builder
	buffer.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")
	rawLen := 64
	keyLen := len(privateKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen
	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(privateKey[start:])
		} else {
			buffer.WriteString(privateKey[start:end])
		}
		buffer.WriteByte('\n')
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END RSA PRIVATE KEY-----\n")
	pKey = buffer.String()
	return
}

// FormatAlipayPublicKey 格式化支付宝普通支付宝公钥
func FormatAlipayPublicKey(publicKey string) (pKey string) {
	var buf strings.Builder
	buf.WriteString("-----BEGIN PUBLIC KEY-----\n")
	rawLen := 64
	keyLen := len(publicKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen
	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buf.WriteString(publicKey[start:])
		} else {
			buf.WriteString(publicKey[start:end])
		}
		buf.WriteByte('\n')
		start += rawLen
		end = start + rawLen
	}
	buf.WriteString("-----END PUBLIC KEY-----\n")
	pKey = buf.String()
	return
}

// RSA解密数据
// t：PKCS1 或 PKCS8
// cipherData：加密字符串byte数组
// privateKey：私钥
func RsaDecrypt(t PKCSType, cipherData []byte, privateKey string) (originData []byte, err error) {
	var (
		key *rsa.PrivateKey
	)

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("privateKey decode error")
	}

	switch t {
	case PKCS1:
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	case PKCS8:
		pkcs8Key, e := x509.ParsePKCS8PrivateKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		pk8, ok := pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("parse PKCS8 key error")
		}
		key = pk8
	default:
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	originBytes, err := rsa.DecryptPKCS1v15(rand.Reader, key, cipherData)
	if err != nil {
		return nil, fmt.Errorf("xrsa.DecryptPKCS1v15：%w", err)
	}
	return originBytes, nil
}

// RSA解密数据
// OAEPWithSHA-256AndMGF1Padding
func RsaDecryptOAEP(h hash.Hash, t PKCSType, privateKey string, ciphertext, label []byte) (originData []byte, err error) {
	var (
		key *rsa.PrivateKey
	)

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("privateKey decode error")
	}

	switch t {
	case PKCS1:
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	case PKCS8:
		pkcs8Key, e := x509.ParsePKCS8PrivateKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		pk8, ok := pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("parse PKCS8 key error")
		}
		key = pk8
	default:
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	originBytes, err := rsa.DecryptOAEP(h, rand.Reader, key, ciphertext, label)
	if err != nil {
		return nil, err
	}
	return originBytes, nil
}

// RSA加密数据
// t：PKCS1 或 PKCS8
// originData：原始字符串byte数组
// publicKey：公钥
func RsaEncrypt(t PKCSType, originData []byte, publicKey string) (cipherData []byte, err error) {
	var (
		key *rsa.PublicKey
	)

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("publicKey decode error")
	}

	switch t {
	case PKCS1:
		pkcs1Key, e := x509.ParsePKCS1PublicKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		key = pkcs1Key
	case PKCS8:
		pkcs8Key, e := x509.ParsePKIXPublicKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		pk8, ok := pkcs8Key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("parse PKCS8 key error")
		}
		key = pk8
	default:
		pkcs1Key, e := x509.ParsePKCS1PublicKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		key = pkcs1Key
	}

	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, key, originData)
	if err != nil {
		return nil, fmt.Errorf("xrsa.EncryptPKCS1v15：%w", err)
	}
	return cipherBytes, nil
}

// RSA加密数据
// OAEPWithSHA-256AndMGF1Padding
func RsaEncryptOAEP(h hash.Hash, t PKCSType, publicKey string, originData, label []byte) (cipherData []byte, err error) {
	var (
		key *rsa.PublicKey
	)
	if len(originData) > 190 {
		return nil, errors.New("message too long for RSA public key size")
	}
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("publicKey decode error")
	}

	switch t {
	case PKCS1:
		pkcs1Key, e := x509.ParsePKCS1PublicKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		key = pkcs1Key
	case PKCS8:
		pkcs8Key, e := x509.ParsePKIXPublicKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		pk8, ok := pkcs8Key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("parse PKCS8 key error")
		}
		key = pk8
	default:
		pkcs1Key, e := x509.ParsePKCS1PublicKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		key = pkcs1Key
	}

	cipherBytes, err := rsa.EncryptOAEP(h, rand.Reader, key, originData, label)
	if err != nil {
		return nil, err
	}
	return cipherBytes, nil
}
