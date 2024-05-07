package compress

import "github.com/golang/snappy"

func SnappyCompress(buf []byte) []byte {
	return snappy.Encode(nil, buf)
}

func SnappyUnCompress(buf []byte) ([]byte, error) {
	return snappy.Decode(nil, buf)
}
