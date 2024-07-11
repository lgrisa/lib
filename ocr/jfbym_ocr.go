package ocr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const CustomUrl = "http://api.jfbym.com/api/YmServer/customApi"
const Token = "utLwhPE31v5BiEMNPtsvy4BEY8B-GRdFyL4047s6EIk"

func CloudCodeCommonVerify(imageBody []byte) (string, error) {

	//# 数英汉字类型
	//# 通用数英1-4位 10110
	//# 通用数英5-8位 10111
	//# 通用数英9~11位 10112
	//# 通用数英12位及以上 10113
	//# 通用数英1~6位plus 10103
	//# 定制-数英5位~qcs 9001
	//# 定制-纯数字4位 193
	//# 中文类型
	//# 通用中文字符1~2位 10114
	//# 通用中文字符 3~5位 10115
	//# 通用中文字符6~8位 10116
	//# 通用中文字符9位及以上 10117
	//# 定制-XX西游苦行中文字符 10107
	//# 计算类型
	//# 通用数字计算题 50100
	//# 通用中文计算题 50101
	//# 定制-计算题 cni 452
	config := map[string]interface{}{}
	config["image"] = base64.StdEncoding.EncodeToString(imageBody)
	config["type"] = "10103"
	config["token"] = Token
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("关闭失败：%v\n", err)
		}
	}(resp.Body)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data), err)
	return string(data), nil
}
