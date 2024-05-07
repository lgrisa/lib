package weight

import (
	"github.com/lgrisa/lib/utils/random"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
)

var (
	ErrChoiceNotEnough = errors.Errorf("weight.len < n")
	ErrZeroWeight      = errors.Errorf("weight is 0")
	ErrLogic           = errors.Errorf("error logic")
)

func RandomN(weight []uint64, n int) ([]int, error) {

	if n <= 0 {
		return nil, nil
	}

	count := len(weight)
	if count <= 0 || count < n {
		return nil, ErrChoiceNotEnough
	}

	if count == n {
		a := make([]int, count)
		for i := 0; i < count; i++ {
			a[i] = i
		}
		return a, nil
	}

	// 循环n次，每次随机一个出来，然后
	var totalWeight uint64
	for _, w := range weight {
		if w <= 0 {
			return nil, ErrZeroWeight
		}
		totalWeight += w
	}

	copyWeight := make([]uint64, count)
	copy(copyWeight, weight)

	index := make([]int, 0, n)

out:
	for i := 0; i < n; i++ {
		x := rand.Uint64() % totalWeight

		// 从第一个开始找起，找到第一个 cw <= x < cw+w
		var cur uint64
		for i := 0; i < count; i++ {

			w := copyWeight[i]

			if w > 0 && x < cur+w {
				totalWeight -= w
				copyWeight[i] = 0

				index = append(index, i)
				continue out
			}

			cur += w
		}

		return nil, ErrLogic
	}

	return index, nil
}

func MustRandomFilterN[T interface{}](weight []T, n int, getW func(w T) uint64) []T {
	ret, err := RandomFilterN(weight, n, getW)
	if err != nil {
		logrus.WithError(err).Errorf("MustRandomFilterN error")

		// 直接随机一下，什么都不管了
		indexArray := random.NewMNIntIndexArray(len(weight), n)
		result := make([]T, 0, len(indexArray))
		for _, idx := range indexArray {
			result = append(result, weight[idx])
		}
		return result
	}

	if len(ret) == 0 {
		indexArray := random.NewMNIntIndexArray(len(weight), n)
		result := make([]T, 0, len(indexArray))
		for _, idx := range indexArray {
			result = append(result, weight[idx])
		}
		return result
	}

	return ret
}

func RandomFilterN[T interface{}](weight []T, n int, getW func(w T) uint64) ([]T, error) {

	var weightArray []uint64
	var dataArray []T
	for _, v := range weight {
		w := getW(v)
		if w > 0 {
			weightArray = append(weightArray, w)
			dataArray = append(dataArray, v)
		}
	}

	if len(dataArray) <= n {
		return dataArray, nil
	}

	indexArray, err := RandomN(weightArray, n)
	if err != nil {
		return nil, err
	}

	result := make([]T, 0, n)
	for _, idx := range indexArray {
		result = append(result, dataArray[idx])
	}
	return result, nil
}
