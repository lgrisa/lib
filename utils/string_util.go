package utils

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// Capitalize 首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// SubString 字符串截取
func SubString(str string, startIndex int, length int) string {
	rs := []rune(str)
	return string(rs[startIndex : startIndex+length])
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

func ParseSize(size string) (int64, string) {
	//默认大小为100MB
	re, _ := regexp.Compile("[0-9]+")
	unit := string(re.ReplaceAll([]byte(size), []byte("")))
	num, _ := strconv.ParseInt(strings.Replace(size, unit, "", 1), 10, 64)
	unit = strings.ToUpper(unit)
	var byteNum int64 = 0
	switch unit {
	case "B":
		byteNum = num
	case "KB":
		byteNum = num * KB
	case "MB":
		byteNum = num * MB
	case "GB":
		byteNum = num * GB
	case "TB":
		byteNum = num * TB
	case "PB":
		byteNum = num * PB
	default:
		num = 0
		byteNum = 0
	}
	if num == 0 {
		log.Println("ParseSize 仅支持B KB MB GB TB PB")
		num = 100
		unit = "MB"
		byteNum = num * MB
	}
	sizeStr := strconv.FormatInt(num, 10) + unit
	return byteNum, sizeStr
}

// Byte2String 快速把一个[]byte转换为string. 使string在背后直接使用这个[]byte作为它的数据而不再copy一份
// 调用之后不能再修改原始的[]byte中的数据, 不然会导致string被修改
func Byte2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Byte 快速把一个string转换为[]byte. 直接使用string背后的[]byte
// 调用之后不能修改[]byte中的数据, 不然会导致string被修改
func String2Byte(b string) []byte {
	return *(*[]byte)(unsafe.Pointer(&b))
}

func GetCharLen(s string) int {

	charLen := 0
	for _, r := range []rune(s) {
		n := utf8.RuneLen(r)
		switch n {
		case -1:
			// 当成最大的unicode字符处理
			charLen += 4
		case 1:
			charLen += 1
		default:
			// 多字节字符，一个当成2个
			charLen += 2
		}
	}

	return charLen
}

func TruncateCharLen(s string, maxCharLen int) string {
	runeArray := []rune(s)
	b := &bytes.Buffer{}

	charLen := 0
	for _, r := range runeArray {
		n := utf8.RuneLen(r)
		switch n {
		case -1:
			// 跳过这种字符
			continue
		case 1:
			charLen += 1
		default:
			// 多字节字符，一个当成2个
			charLen += 2
		}

		if charLen > maxCharLen {
			break
		}
		b.WriteRune(r)
	}

	return b.String()
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

func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func Split2(s, sep, sep2 string) [][]string {
	array := strings.Split(s, sep)
	rets := make([][]string, 0, len(array))
	for _, v := range array {
		rets = append(rets, Split(v, sep2))
	}
	return rets
}

func SplitToNumber(s, sep string) ([]int64, error) {
	var ret []int64
	for _, v := range strings.Split(s, sep) {
		i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err == nil {
			ret = append(ret, i)
		} else {
			return nil, errors.Wrapf(err, "解析数字失败，字符串：%s", v)
		}
	}
	return ret, nil
}
func SplitToF64(s, sep string) ([]float64, error) {
	var ret []float64
	for _, v := range strings.Split(s, sep) {
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err == nil {
			ret = append(ret, f)
		} else {
			return nil, errors.Wrapf(err, "解析数字失败，字符串：%s", v)
		}
	}
	return ret, nil
}

func SplitToNumber2(s, sep, sep2 string) ([][]int64, error) {
	array := strings.Split(s, sep)

	rets := make([][]int64, 0, len(array))
	for _, v := range array {
		ret, err := SplitToNumber(v, sep2)
		if err != nil {
			return nil, errors.Wrapf(err, s)
		}
		rets = append(rets, ret)
	}
	return rets, nil
}

func SplitToI64F64Map(s, sep, sep2 string) (map[int64]float64, error) {
	array := strings.Split(s, sep)

	rets := make(map[int64]float64, len(array))
	for _, v := range array {
		array := strings.Split(v, sep2)
		if len(array) != 2 {
			return nil, errors.New("解析失败, len(array) != 2")
		}

		i, err := strconv.ParseInt(strings.TrimSpace(array[0]), 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "解析数字失败，字符串：%s", array[0])
		}

		f, err := strconv.ParseFloat(strings.TrimSpace(array[1]), 64)
		if err != nil {
			return nil, errors.Wrapf(err, "解析数字失败，字符串：%s", array[1])
		}

		if _, exist := rets[i]; exist {
			return nil, errors.New("解析失败, 存在重复的key")
		}

		rets[i] = f
	}
	return rets, nil
}

func SplitToStrMap(s, sep, sep2 string) (map[string]string, error) {
	array := strings.Split(s, sep)

	rets := make(map[string]string, len(array))
	for _, v := range array {
		array := strings.Split(v, sep2)
		if len(array) != 2 {
			return nil, errors.Errorf("解析失败, len(array) != 2, %v", s)
		}
		key := strings.TrimSpace(array[0])
		value := strings.TrimSpace(array[1])

		if _, exist := rets[key]; exist {
			return nil, errors.Errorf("解析失败, duplicate key, %v", s)
		}

		rets[key] = value
	}
	return rets, nil
}
