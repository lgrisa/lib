package utils

import (
	"bytes"
	"unicode/utf8"
)

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
