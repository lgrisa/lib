package utils

import (
	"archive/zip"
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ZipPackDir(root string) ([]byte, error) {

	fileMap := make(map[string][]byte)
	handleFile := func(root, zipPrefix string, f os.DirEntry) error {
		if f.IsDir() {
			return nil
		}

		if strings.HasPrefix(f.Name(), "~$") {
			return nil
		}

		filename := filepath.Join(root, f.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			return errors.Wrapf(err, "读取文件失败，file: %s", filename)
		}

		fileMap[zipPrefix+f.Name()] = data
		return nil
	}

	fs, err := doReadDir(root)
	if err != nil {
		return nil, errors.Wrapf(err, "读取文件夹失败，root: %s", root)
	}
	for _, f := range fs {
		if err = handleFile(root, "", f); err != nil {
			return nil, errors.Errorf("处理文件失败，root: %s, file: %s err: %v", root, f.Name(), err)
		}
	}

	zipFile, err := PackZipData(fileMap)
	if err != nil {
		return nil, errors.Wrap(err, "压缩文件失败")
	}

	return zipFile, nil
}

func doReadDir(dir string) ([]os.DirEntry, error) {
	if strings.HasSuffix(dir, "/") {
		dir = dir[:len(dir)-1]
	}

	return os.ReadDir(dir)
}

func PackZipData(fileBytes map[string][]byte) ([]byte, error) {

	var ks []string
	for k := range fileBytes {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	// 压缩成zip
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)

	for _, k := range ks {
		f, err := w.Create(k)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(fileBytes[k])
		if err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func UnpackZipData(data []byte) (map[string][]byte, error) {
	buf := bytes.NewReader(data)
	r, err := zip.NewReader(buf, int64(len(data)))
	if err != nil {
		return nil, err
	}

	fileBytes := make(map[string][]byte)
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		b, err := io.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		filename := strings.ReplaceAll(f.Name, "\\", "/")
		fileBytes[filename] = b
	}

	return fileBytes, nil
}
