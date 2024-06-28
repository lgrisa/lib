package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// Capitalize 首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		// 空格后字母大写
		isBlank := i > 0 && vv[i-1] == ' '

		if i == 0 || isBlank {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
			}
		}

		upperStr += string(vv[i])
	}
	return upperStr
}

// ToStringE interface to string
func ToStringE(i interface{}) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatInt(int64(s), 10), nil
	case uint64:
		return strconv.FormatInt(int64(s), 10), nil
	case uint32:
		return strconv.FormatInt(int64(s), 10), nil
	case uint16:
		return strconv.FormatInt(int64(s), 10), nil
	case uint8:
		return strconv.FormatInt(int64(s), 10), nil
	case []byte:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
)

// GetByte2StringPointer 快速把一个[]byte转换为string. 使string在背后直接使用这个[]byte作为它的数据而不再copy一份
// 调用之后不能再修改原始的[]byte中的数据, 不然会导致string被修改
func GetByte2StringPointer(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// GetString2BytePointer 快速把一个string转换为[]byte. 直接使用string背后的[]byte
// 调用之后不能修改[]byte中的数据, 不然会导致string被修改
func GetString2BytePointer(b string) []byte {
	return *(*[]byte)(unsafe.Pointer(&b))
}

// ReplaceInvalidChar 替换非法的字符
func ReplaceInvalidChar(s string) string {
	runeArray := []rune(s)
	buffer := bytes.NewBuffer(make([]byte, 0, len(runeArray)))

	for _, r := range runeArray {
		if isValidRune(r, true) {
			buffer.WriteRune(r)
		}
	}

	return strings.TrimSpace(buffer.String())
}

func IsValidName(s string) bool {

	// 名字，不能空格，符号开头
	// 不能包含非法字符

	runeArray := []rune(s)
	lastIdx := len(runeArray) - 1
	for i, r := range runeArray {
		puntValid := !(i == 0 || i == lastIdx)

		// 开头必须是字母，或者数字，不能是符号开头
		if !isValidRune(r, puntValid) {
			return false
		}
	}

	return true
}

func HaveInvalidChar(s string) bool {
	runeArray := []rune(s)

	for _, r := range runeArray {
		if !isValidRune(r, true) {
			return true
		}
	}

	return false
}

func isValidRune(r rune, punctValid bool) bool {
	switch r {
	case '_', ' ':
		if !punctValid {
			return false
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':

	default:
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
