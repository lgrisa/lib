package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Hash(plainText string) string {
	var h = md5.New()
	h.Write([]byte(plainText))
	return hex.EncodeToString(h.Sum(nil))
}
