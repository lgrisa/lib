package util

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/golang/snappy"
	"github.com/lgrisa/lib/utils/pbutil"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

var (
	pool = pbutil.Pool
)

type Marshaler interface {
	Marshal() (dAtA []byte, err error)
}

func SafeMarshal(m Marshaler) []byte {
	data, err := m.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("safe.Marshal fail")
	}
	return data
}

type proto interface {
	Size() int
	MarshalTo([]byte) (int, error)
	MarshalToSizedBuffer([]byte) (int, error)
}

func NewProtoMsg(object proto, head []byte, msgName string, static bool) pbutil.Buffer {
	headLen := len(head)
	_size := headLen + object.Size()
	result, buf := allocBuff(_size, static)

	copy(buf, head)
	if _, err := object.MarshalToSizedBuffer(buf[headLen:]); err != nil {
		result.Free()
		logrus.WithError(err).Errorf("%s.Marshal fail", msgName)
		return pbutil.Empty
	} else {
		return result
	}
}

const (
	compress_module_id          = 0
	snappy_compress_sequence_id = 0
	gzip_compress_sequence_id   = 1
)

func NewCompressMsg(object proto, head []byte, msgName string, static bool) pbutil.Buffer {
	return NewSnappyCompressMsg(object, head, msgName, static)
}

var s2c_snappy_compress_msg = [...]byte{compress_module_id, snappy_compress_sequence_id}

func NewSnappyCompressMsg(object proto, head []byte, msgName string, static bool) pbutil.Buffer {
	return newCompressMsg(object, head, msgName, static, s2c_snappy_compress_msg[:], func(buf []byte) []byte {
		return snappy.Encode(nil, buf)
	})
}

var s2c_gzip_compress_msg = [...]byte{compress_module_id, gzip_compress_sequence_id}

func NewGzipCompressMsg(object proto, head []byte, msgName string, static bool) pbutil.Buffer {
	return newCompressMsg(object, head, msgName, static, s2c_gzip_compress_msg[:], func(buf []byte) []byte {
		var b bytes.Buffer
		w, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)
		w.Write(buf)
		w.Close()
		return b.Bytes()
	})
}

func newCompressMsg(object proto, head []byte, msgName string, static bool, compressHead []byte, compressFunc func(buf []byte) []byte) pbutil.Buffer {

	headLen := len(head)

	buf := pool.Alloc(headLen + object.Size())
	defer buf.Free()

	copy(buf, head)
	if _, err := object.MarshalToSizedBuffer(buf[headLen:]); err != nil {
		logrus.WithError(err).Errorf("%s.Marshal fail", msgName)
		return pbutil.Empty
	} else {
		return NewBytesMsg(compressFunc(buf), compressHead, static)
	}
}

func NewBytesMsg(data, head []byte, static bool) pbutil.Buffer {
	headLen := len(head)
	_size := headLen + len(data)
	result, buf := allocBuff(_size, static)

	copy(buf, head)
	copy(buf[headLen:], data)
	return result
}

func IsCompressMsg(moduleID, sequenceID int) bool {
	if compress_module_id == moduleID {
		switch sequenceID {
		case snappy_compress_sequence_id, gzip_compress_sequence_id:
			return true
		}
	}
	return false
}

var errUnkownCompressSequenceId = errors.New("Unkown compress msg sequence")

func UncompressMsg(sequenceID int, data []byte) ([]byte, error) {
	switch sequenceID {
	case snappy_compress_sequence_id:
		return snappy.Decode(nil, data)
	case gzip_compress_sequence_id:
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer r.Close()
		uncomress, err := ioutil.ReadAll(r)
		return uncomress, err
	default:
		return nil, errUnkownCompressSequenceId
	}
}

func poolAlloc(size int, static bool) pbutil.Buffer {
	if static {
		return pbutil.StaticBuffer(make([]byte, size))
	} else {
		return pbutil.NewRecycleBuffer(size)
	}
}

func allocBuff(_size int, static bool) (result pbutil.Buffer, buf []byte) {
	switch {
	case _size <= 127:
		result = poolAlloc(_size+2, static)
		rb := result.Buffer()
		rb[0] = 0
		buf = rb[1:]
		buf[0] = uint8(_size)
		buf = buf[1:]
	case _size <= 16383:
		result = poolAlloc(_size+3, static)
		rb := result.Buffer()
		rb[0] = 0
		buf = rb[1:]
		buf[0], buf[1] = (0x80 | uint8(_size&0x7f)), uint8(_size>>7)
		buf = buf[2:]
	default:
		result = poolAlloc(_size+4, static)
		rb := result.Buffer()
		rb[0] = 0
		buf = rb[1:]
		buf[0], buf[1], buf[2] = (0x80 | uint8(_size&0x7f)), (0x80 | uint8((_size>>7)&0x7f)), uint8(_size>>14)
		buf = buf[3:]
	}

	return
}
