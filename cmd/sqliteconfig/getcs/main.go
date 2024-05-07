package main

import (
	"bytes"
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lgrisa/lib/utils"
	"github.com/lgrisa/lib/utils/compress"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

var httpAddr = flag.String("httpAddr", "", "获取CS http地址")
var confPath = flag.String("confPath", "conf/", "配置表路径")
var csPath = flag.String("csPath", "", "cs表路径")
var protoPath = flag.String("protoPath", "", "proto路径")

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	// code == 200时，返回的数据
	Data *ConfigCsJson `json:"data"`
}

type ConfigCsJson struct {
	FileMd5 string `json:"file_md5"`

	FileSize int64 `json:"file_size"`

	Body []byte `json:"cs_body"`
}

func main() {
	flag.Parse()

	// 遍历conf文件夹
	data, err := compress.ZipPackDir(*confPath)
	if err != nil {
		fmt.Println("遍历文件夹压缩生成Data失败：", err)
		return
	}

	body, err := doHttpGet(*httpAddr, data)

	if err != nil {
		fmt.Println("http请求失败：", err)
		return
	}

	resp := &response{}
	err = jsoniter.Unmarshal(body, resp)
	if err != nil {
		fmt.Println("解析json失败：", err)
		return
	}

	fmt.Printf("FileMd5: %s, FileSize: %d\n", resp.Data.FileMd5, resp.Data.FileSize)

	fileMap, err := compress.UnpackZipData(resp.Data.Body)
	if err != nil {
		fmt.Println("解压zip失败：", err)
		return
	}

	//清除原来的Proto文件夹直接重新生成
	_, errCsIsExist := os.Stat(*csPath)

	if !os.IsNotExist(errCsIsExist) {
		fmt.Println("删除原来的cs文件夹")
		_ = os.RemoveAll(*csPath)
	}

	if *csPath != "" {
		for filename, fileBytes := range fileMap {
			if !strings.HasSuffix(filename, ".cs") {
				continue
			}

			// 将cs写入到本地磁盘
			targetPath := *csPath + filename
			if err = os.WriteFile(targetPath, fileBytes, os.ModePerm); err != nil {
				fmt.Printf("写入 cs: %s 失败：%s\n", targetPath, err)
				return
			}
		}
	}

	if *protoPath != "" {
		for filename, fileBytes := range fileMap {
			if !strings.HasSuffix(filename, ".proto") {
				continue
			}

			// 将proto写入到本地磁盘
			targetPath := *protoPath + filename
			if err = os.WriteFile(targetPath, fileBytes, os.ModePerm); err != nil {
				fmt.Printf("写入Proto: %s 失败：%s\n", targetPath, err)
				return
			}
		}
	}
}

func doHttpGet(httpAddr string, data []byte) ([]byte, error) {
	md5String := utils.Md5String(data)
	size := len(data)

	// 发送到服务器进行解析
	r, err := http.NewRequest(http.MethodGet, httpAddr, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrapf(err, "new request fail")
	}

	r.Header.Set("file_md5", md5String)
	r.Header.Set("file_size", fmt.Sprintf("%d", size))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.Wrapf(err, "do request fail")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read response fail")
	}

	logrus.WithField("code", resp.StatusCode).Debug(string(body))

	return body, nil
}
