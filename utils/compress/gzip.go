package compress

import (
	"bytes"
	"compress/gzip"
	"io"
)

func GzipCompress(buf []byte) []byte {
	return gzipCompressionWithLevel(buf, gzip.DefaultCompression)
}

func GzipCompress2(buf1, buf2 []byte) []byte {
	return gzipCompressionArrayWithLevel(gzip.DefaultCompression, buf1, buf2)
}

func GzipUnCompress(buf []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	return io.ReadAll(r)
}

func gzipCompressionWithLevel(buf []byte, level int) []byte {
	var b bytes.Buffer
	w, _ := gzip.NewWriterLevel(&b, level)
	_, _ = w.Write(buf)
	_ = w.Close()
	return b.Bytes()
}

func gzipCompressionArrayWithLevel(level int, buffs ...[]byte) []byte {
	var b bytes.Buffer
	w, _ := gzip.NewWriterLevel(&b, level)
	for _, buf := range buffs {
		_, _ = w.Write(buf)
	}
	_ = w.Close()
	return b.Bytes()
}
