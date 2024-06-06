package config

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/lgrisa/lib/utils/file"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func GetConfigVersion() (string, error) {
	// 读取配置的文件夹，每个文件都加载上来，计算

	s, err := file.FindConfigPath("conf")
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
