package file

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
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

func GetConfigVersion() (string, error) {
	// 读取配置的文件夹，每个文件都加载上来，计算

	s, err := FindConfigPath("conf")
	if err != nil {
		return "", errors.Errorf("生成配置版本号，获取配置路径失败")
	}

	hash := md5.New()
	err = filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			// 隐藏文件，跳过
			return nil
		}

		if strings.Contains(path, ".svn") {
			// 跳过svn内容
			return nil
		}

		if strings.Contains(path, ".git") {
			// 跳过git内容
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if len(b) > 0 {
			hash.Write(b)
		}
		return nil
	})

	if err != nil {
		return "", errors.Errorf("生成配置版本号，遍历配置文件夹失败")
	}

	sum := hash.Sum(nil)
	for i := 1; i < len(sum)/2; i++ {
		sum[0] ^= sum[i*2]
		sum[1] ^= sum[i*2+1]
	}

	return hex.EncodeToString(sum[:2]), nil
}
