package utils

import (
	//	"fmt"
	"bytes"
	"github.com/lgrisa/lib/utils/sortkeys"
	"sort"
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

func Int32Duplicate(array []int32) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Int32DuplicateIgnoreZero(array []int32) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		if x == 0 {
			continue
		}

		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Int32CountIgnoreZero(array []int32) int {
	c := 0
	for _, v := range array {
		if v != 0 {
			c++
		}
	}

	return c
}

func Int32AnyZero(array []int32) bool {
	for _, v := range array {
		if v == 0 {
			return true
		}
	}

	return false
}

func Int32AnyLe0(array []int32) bool {
	for _, v := range array {
		if v <= 0 {
			return true
		}
	}

	return false
}

func Int32AnyLt0(array []int32) bool {
	for _, v := range array {
		if v < 0 {
			return true
		}
	}

	return false
}

func Int32AnyGt0(array []int32) bool {
	for _, v := range array {
		if v > 0 {
			return true
		}
	}

	return false
}

func IntDuplicate(array []int) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Int64Duplicate(array []int64) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Int64AnyGtN(array []int64, n int64) bool {
	for _, v := range array {
		if v > n {
			return true
		}
	}
	return false
}

func Int64AnyGeN(array []int64, n int64) bool {
	for _, v := range array {
		if v >= n {
			return true
		}
	}
	return false
}

func Uint64Duplicate(array []uint64) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func StringDuplicate(array []string) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func StringNilOrDuplicate(array []string) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		if len(x) <= 0 {
			return true
		}

		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Uint64Sorted(array []uint64) bool {
	return sort.IsSorted(sortkeys.Uint64Slice(array))
}

func Float64AnyLe0(array []float64) bool {
	for _, v := range array {
		if v <= 0 {
			return true
		}
	}

	return false
}

func Float64AnyLt0(array []float64) bool {
	for _, v := range array {
		if v < 0 {
			return true
		}
	}

	return false
}
