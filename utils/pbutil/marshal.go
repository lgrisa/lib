package pbutil

import (
	"github.com/lgrisa/lib/utils/pool"
)

type Marshaller interface {
	MarshalToSizedBuffer([]byte) (int, error)
	Size() int
}

func Marshal(v Marshaller) (pool.Buffer, error) {
	size := v.Size()
	buf := pool.Pool.Alloc(size)

	if _, err := v.MarshalToSizedBuffer(buf); err != nil {
		buf.Free()
		return nil, err
	}

	return buf, nil
}
