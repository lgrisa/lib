package varint

import (
	"encoding/binary"
	"math/bits"
)

// ensure 'sz' extra bytes in 'b' btw len(b) and cap(b)
func ensure(b []byte, sz int) ([]byte, int) {
	l := len(b)
	c := cap(b)
	if c-l < sz {
		o := make([]byte, (2*c)+sz) // exponential growth
		n := copy(o, b)
		return o[:n+sz], n
	}
	return b[:l+sz], l
}

// UvarintSize returns the size (in bytes) of `num` encoded as a unsigned varint.
//
// This may return a size greater than MaxUvarintLen63, which would be an
// illegal value, and would be rejected by readers.
func UvarintSize(num uint64) int {
	bits := bits.Len64(num)
	q, r := bits/7, bits%7
	size := q
	if r > 0 || size == 0 {
		size++
	}
	return size
}

// AppendInt64 appends an int64 to the slice
func AppendUint64(b []byte, u uint64) []byte {
	o, n := ensure(b, UvarintSize(u))
	binary.PutUvarint(o[n:], u)
	return o
}

func AppendUint32(b []byte, u uint32) []byte {
	return AppendUint64(b, uint64(u))
}

func AppendUint16(b []byte, u uint16) []byte {
	return AppendUint64(b, uint64(u))
}

func AppendUint8(b []byte, u uint8) []byte {
	return AppendUint64(b, uint64(u))
}

// AppendInt64 appends an int64 to the slice
func AppendInt64(b []byte, x int64) []byte {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return AppendUint64(b, ux)
}

func AppendInt32(b []byte, x int32) []byte {
	return AppendInt64(b, int64(x))
}

func AppendInt16(b []byte, x int16) []byte {
	return AppendInt64(b, int64(x))
}

func AppendInt8(b []byte, x int8) []byte {
	return AppendInt64(b, int64(x))
}
