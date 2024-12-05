package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lgrisa/lib/tool/sqliteconfig/mgr"
	"github.com/lgrisa/lib/utils"
	"github.com/lgrisa/lib/utils/compress"
	"net/http"
	"os"
	"strings"
)

var httpAddr = flag.String("httpAddr", "", "http sever地址")
var confPath = flag.String("confPath", "conf/", "conf文件夹路径")
var csPath = flag.String("csPath", "", "cs写入路径")
var protoPath = flag.String("protoPath", "", "proto写入路径")

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

	data, err := compress.ZipPackDir(*confPath)
	if err != nil {
		utils.LogErrorF("compress.ZipPackDir fail: %v", err)
		return
	}

	_, _, body, err := utils.Request(context.Background(), *httpAddr, http.MethodGet, map[string]string{
		mgr.ReqHeaderFileMd5:  utils.Md5String(data),
		mgr.ReqHeaderFileSize: fmt.Sprintf("%d", len(data)),
	}, bytes.NewReader(data))

	if err != nil {
		utils.LogErrorF("request fail: %v", err)
		return
	}

	resp := &response{}
	if err = jsoniter.Unmarshal(body, resp); err != nil {
		utils.LogErrorF("json.Unmarshal fail: %v", err)
		return
	}

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
				utils.LogErrorF("write cs fail: %v", err)
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
				utils.LogErrorF("write proto fail: %v", err)
				return
			}
		}
	}
}
