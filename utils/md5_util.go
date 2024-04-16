package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func Md5String(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

func Md5PathString(path string) (string, error) {
	b, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:]), nil
}

func ReadFileMd5(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return Md5String(data), nil
}
