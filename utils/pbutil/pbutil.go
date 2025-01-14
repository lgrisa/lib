package pbutil

import (
	"github.com/lgrisa/lib/utils/pool"
	"sync/atomic"
)

var (
	Pool = pool.Pool

	Empty = StaticBuffer([]byte{})

	freeCacheSize int

	duplicateFreeCallback func([]byte)
)

func InitDuplicateFreeCallback(f func([]byte), toSet int) {
	if toSet < 0 {
		toSet = 0
	}

	duplicateFreeCallback = f
	freeCacheSize = toSet
}

type Buffer interface {
	Buffer() []byte
	Free() bool
}

var _ Buffer = (StaticBuffer)(nil)

type StaticBuffer []byte

func (d StaticBuffer) Buffer() []byte {
	return d
}

func (d StaticBuffer) Free() bool {
	return true
}

func NewRecycleBuffer(size int) Buffer {
	return newRecycleBuffer(Pool.Alloc(size))
}

var _ Buffer = (*recycleBuffer)(nil)

func newRecycleBuffer(buf []byte) *recycleBuffer {
	return &recycleBuffer{buf: buf}
}

type recycleBuffer struct {
	buf []byte

	free atomic.Bool

	// 缓存消息前7个字节，用于查问题时知道哪个消息重复Free
	freeCache []byte
}

func (d *recycleBuffer) Buffer() []byte {
	return d.buf
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *recycleBuffer) Free() bool {
	if d.free.CompareAndSwap(false, true) {
		if freeCacheSize > 0 {
			// 缓存前n个字节
			d.freeCache = make([]byte, min(len(d.buf), freeCacheSize))
			copy(d.freeCache, d.buf)
		}
		Pool.Free(d.buf)
		d.buf = nil
		return true
	}

	if duplicateFreeCallback != nil {
		duplicateFreeCallback(d.freeCache)
	}

	return false
}
