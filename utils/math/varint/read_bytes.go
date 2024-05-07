package varint

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"math"
)

var BufTooSmall = errors.Errorf("binary: varint buf too small")
var Overflow = errors.Errorf("binary: varint overflows a 64-bit integer")

// ReadUint64Bytes tries to read a uint64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - BufTooSmall (too few bytes)
// - Overflow
func ReadUint64Bytes(b []byte) (u uint64, o []byte, err error) {
	v, n := binary.Uvarint(b)
	switch {
	case n == 0:
		// buf too small
		return 0, b, BufTooSmall
	case n < 0:
		// value larger than 64 bits (overflow) and -n is the number of bytes read
		return 0, b, Overflow
	}
	return v, b[n:], nil
}

func ReadUint32Bytes(b []byte) (uint32, []byte, error) {
	i, o, err := ReadUint64Bytes(b)
	if i > math.MaxUint32 {
		return 0, o, errors.Errorf("binary: varuint overflows a %v-bit uinteger(%v)", 32, i)
	}
	return uint32(i), o, err
}

func ReadUint16Bytes(b []byte) (uint16, []byte, error) {
	i, o, err := ReadUint64Bytes(b)
	if i > math.MaxUint16 {
		return 0, o, errors.Errorf("binary: varuint overflows a %v-bit uinteger(%v)", 16, i)
	}
	return uint16(i), o, err
}

func ReadUint8Bytes(b []byte) (uint8, []byte, error) {
	i, o, err := ReadUint64Bytes(b)
	if i > math.MaxUint8 {
		return 0, o, errors.Errorf("binary: varuint overflows a %v-bit uinteger(%v)", 8, i)
	}
	return uint8(i), o, err
}

// ReadInt64Bytes tries to read an int64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - BufTooSmall (too few bytes)
// - Overflow
func ReadInt64Bytes(b []byte) (i int64, o []byte, err error) {
	ux, b, err := ReadUint64Bytes(b)
	if err != nil {
		return 0, b, err
	}
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, b, nil
}

func ReadInt32Bytes(b []byte) (int32, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0, o, errors.Errorf("binary: varint overflows a %v-bit integer(%v)", 32, i)
	}
	return int32(i), o, err
}

func ReadInt16Bytes(b []byte) (int16, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt16 || i < math.MinInt16 {
		return 0, o, errors.Errorf("binary: varint overflows a %v-bit integer(%v)", 16, i)
	}
	return int16(i), o, err
}

func ReadInt8Bytes(b []byte) (int8, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt8 || i < math.MinInt8 {
		return 0, o, errors.Errorf("binary: varint overflows a %v-bit integer(%v)", 8, i)
	}
	return int8(i), o, err
}
