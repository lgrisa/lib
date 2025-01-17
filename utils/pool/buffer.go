package pool

const (
	MaxCacheSize = 1 << 24 // 16MB
)

var (
	Pool = newSyncPool(16, MaxCacheSize, 2)

	Empty = Buffer(nil)
)

type Buffer []byte

func (r Buffer) Free() Buffer {
	Pool.Free(r)
	return Empty
}
