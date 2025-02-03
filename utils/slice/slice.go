package slice

import "golang.org/x/exp/constraints"

func Contains[T comparable](a []T, b T) bool {
	for _, e := range a {
		if e == b {
			return true
		}
	}
	return false
}

func Remove[T comparable](a []T, b T) []T { // 会改变顺序
	for i := 0; i < len(a); i++ {
		if a[i] == b {
			j := len(a) - 1
			a[i], a[j] = a[j], a[i] // 调换位置而不是直接覆盖，保险点
			a = a[:j]               // 去掉最后位置上的元素
			i--                     // 继续从该位置检查
		}
	}
	return a
}

func RemoveIndex[T comparable](a []T, i int) []T { // 会改变顺序

	j := len(a) - 1
	if i >= 0 && i <= j {
		a[i], a[j] = a[j], a[i] // 调换位置而不是直接覆盖，保险点
		a = a[:j]               // 去掉最后位置上的元素
	}
	return a
}

func Max[T constraints.Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	maxS := slice[0]
	for _, v := range slice {
		if v > maxS {
			maxS = v
		}
	}
	return maxS
}
