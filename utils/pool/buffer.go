package pool

const (
	MAX_CACHE_SIZE = 1 << 24 // 16MB
)

var (
	Pool = newSyncPool(16, MAX_CACHE_SIZE, 2)

	Empty = Buffer(nil)
)

type Buffer []byte

func (r Buffer) Free() Buffer {
	Pool.Free(r)
	return Empty
}
