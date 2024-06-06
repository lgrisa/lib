package file

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetFilePrefix(fileName string) string {
	fileNameAll := path.Base(fileName)
	fileSuffix := path.Ext(fileNameAll)
	filePrefix := fileNameAll[0 : len(fileNameAll)-len(fileSuffix)]

	return filePrefix
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ReplaceFileContentField(filePath string, replaceField string, insteadFile string) error {

	in, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(filePath+".mdf", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("Open write file fail:", err)
		os.Exit(-1)
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	br := bufio.NewReader(in)
	index := 1
	for {
		line, prefix, errRead := br.ReadLine()

		if errRead == io.EOF {
			break
		}

		if errRead != nil {
			return errRead
		}

		if prefix {
			return fmt.Errorf("line too long")
		}

		newLine := strings.Replace(string(line), replaceField, insteadFile, -1)
		_, err = out.WriteString(newLine + "\n")
		if err != nil {
			return err
		}

		index++
	}

	_ = in.Close()

	_ = out.Close()

	//删除原文件
	if err = os.Remove(filePath); err != nil {
		return err
	}

	//重命名新文件
	if err = os.Rename(filePath+".mdf", filePath); err != nil {
		return err
	}

	return nil
}

func FindConfigPath(folderName string) (string, error) {

	path0, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	path0 = strings.ReplaceAll(path0, "\\", "/")

	// 防御性，最多100次
	for i := 0; i < 100; i++ {
		confDir := path.Join(path0, folderName)
		if IsDirExist(confDir) {
			// 文件夹存在
			return confDir, nil
		}

		parent := path.Dir(path0)
		if parent == path0 {
			return "", errors.Errorf("配置文件夹 %s 没找到", folderName)
		}
		path0 = parent
		path0 = strings.ReplaceAll(path0, "\\", "/")
	}

	return "", errors.Errorf("配置文件夹 %s 没找到（100次都找不到）", folderName)
}

func IsDirExist(path string) bool {
	fs, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			logrus.WithError(err).Errorf("os.Data(%s) 出错", path)
		}

		return false
	}

	return fs.IsDir()
}
