package pbutil

import (
	"github.com/lgrisa/lib/cmd/sqliteconfig/mgr/pool"
)

type Marshaler interface {
	MarshalToSizedBuffer([]byte) (int, error)
	Size() int
}

func Marshal(v Marshaler) (pool.Buffer, error) {
	size := v.Size()
	buf := pool.Pool.Alloc(size)

	if _, err := v.MarshalToSizedBuffer(buf); err != nil {
		buf.Free()
		return nil, err
	}

	return buf, nil
}
