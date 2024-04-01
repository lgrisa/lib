package utils

import (
	"math/rand"
	"reflect"
	"time"
	"unsafe"
)

// RandomString 随机生成字符串
func RandomString(l int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	bytes := []byte(str)
	var result = make([]byte, 0, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return BytesToString(result)
}

// BytesToString 0 拷贝转换 slice byte 为 string
func BytesToString(b []byte) (s string) {
	_bpt := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	_spur := (*reflect.StringHeader)(unsafe.Pointer(&s))
	_spur.Data = _bpt.Data
	_spur.Len = _bpt.Len
	return s
}
