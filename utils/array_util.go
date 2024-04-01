package utils

import (
	//	"fmt"
	"bytes"
	"strings"
)

// InArray 传入InterfaceArray 暂时用处不多
func InArray(need interface{}, needArr []interface{}) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

// InStringArray 忽略空格去对比
func InStringArray(value string, targetStringSlice []string, ignoreSpace bool) bool {
	for _, v := range targetStringSlice {
		if ignoreSpace {
			value = strings.TrimSpace(value)
			v = strings.TrimSpace(v)
		}

		if value == v {
			return true
		}

	}
	return false
}

func InUintArray(need uint, needArr []uint) bool {

	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

// Intersect 取交集
func Intersect(slice1 []string, slice2 []string) []string {
	var diffSlice []string
	for _, v := range slice1 {
		for _, v2 := range slice2 {
			if v == v2 {
				diffSlice = append(diffSlice, v)
				break
			}
		}
	}

	return diffSlice
}

// BytesCombine 合并byteArray
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
