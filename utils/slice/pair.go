package slice

import "github.com/lgrisa/lib/utils/math/imath"

func NewKvs[K, V interface{}](ks []K, vs []V) []*Kv[K, V] {
	n := imath.Min(len(ks), len(vs))
	kvs := make([]*Kv[K, V], 0, n)
	for i := 0; i < n; i++ {
		kvs = append(kvs, &Kv[K, V]{
			Key:   ks[i],
			Value: vs[i],
		})
	}
	return kvs
}

type Kv[K, V interface{}] struct {
	Key   K
	Value V
}
