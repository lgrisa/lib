package utils

import (
	"github.com/lgrisa/lib/utils/math/imath"
	"github.com/lgrisa/lib/utils/math/u64"
	"github.com/lgrisa/lib/utils/random"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sort"
)

var (
	ErrChoiceNotEnough = errors.Errorf("weight.len < n")
	ErrZeroWeight      = errors.Errorf("weight is 0")
	ErrLogic           = errors.Errorf("error logic")
)

type WeightRandomE struct {
	totalWeight int
	maxIndex    int
	n           int
	rankArray   []uint64
}

func NewWeightRandomE(weight []uint64) (*WeightRandomE, error) {

	if len(weight) == 0 {
		return nil, errors.Errorf("创建权重随机器，权重列表长度为0")
	}

	array := make([]uint64, len(weight))
	maxIndex := len(array) - 1

	var totalWeight uint64
	for i, x := range weight {
		if x == 0 {
			return nil, errors.Errorf("创建权重随机器失败，权重值不能为0，rankArray[%v] == 0", i)
		}

		array[maxIndex-i] = totalWeight // reverse
		totalWeight += x
	}

	r := &WeightRandomE{
		totalWeight: u64.Int(totalWeight),
		maxIndex:    maxIndex,
		n:           len(array),
		rankArray:   array,
	}

	return r, nil
}

func NewRankRandomE(rankArray []uint64) (*WeightRandomE, error) {

	if len(rankArray) == 0 {
		return nil, errors.Errorf("创建权重随机器，rank列表长度为0")
	}

	array := make([]uint64, len(rankArray))
	maxIndex := len(array) - 1

	var prev uint64
	for i, x := range rankArray {
		if x <= prev {
			return nil, errors.Errorf("创建权重随机器，rank列表必须是从小到大的排列顺序")
		}

		array[maxIndex-i] = x
		prev = x
	}

	r := &WeightRandomE{
		totalWeight: u64.Int(array[0]),
		maxIndex:    maxIndex,
		n:           len(array),
		rankArray:   array,
	}

	return r, nil
}

func (r *WeightRandomE) RandomIndex() int {
	return r.Index(uint64(rand.Intn(r.totalWeight)))
}

func (r *WeightRandomE) Index(x uint64) int {
	return imath.Max(0, r.maxIndex-sort.Search(r.n, func(i int) bool { return r.rankArray[i] <= x }))
}

type weight interface {
	GetWeight() uint64
}

func RandomI64Weight(array []int64) int {
	idx, _ := randomFilter(array, func(u int64) uint64 {
		if u > 0 {
			return uint64(u)
		}
		return 0
	})
	return idx
}

func RandomU64Weight(array []uint64) int {
	idx, _ := randomFilter(array, func(u uint64) uint64 {
		return u
	})
	return idx
}

func RandomWeights[T weight](array []T) (t T) {
	return RandomFilter(array, func(t T) uint64 {
		return t.GetWeight()
	})
}

func RandomFilter[T interface{}](array []T, getW func(w T) uint64) (t T) {
	_, t = randomFilter(array, getW)
	return t
}

func randomFilter[T interface{}](array []T, getW func(w T) uint64) (idx int, t T) {

	if len(array) <= 0 {
		return
	}

	totalWeight := uint64(0)
	var weightArray []uint64
	for _, v := range array {
		w := getW(v)
		totalWeight += w
		weightArray = append(weightArray, w)
	}

	if totalWeight > 0 {
		value := uint64(rand.Intn(int(totalWeight)))
		totalWeight = uint64(0)
		for i, v := range array {
			totalWeight += weightArray[i]
			if value < totalWeight {
				return i, v
			}
		}
	}

	idx = rand.Intn(len(array))
	return idx, array[idx]
}

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
