package main

import (
	"context"
	"fmt"
	"github.com/lgrisa/lib/ali/ocr"
	"github.com/lgrisa/lib/utils"
	"os"
	"strconv"
	"time"
)

func main() {
	timeNow := time.Now().UnixMicro()

	timeStr := strconv.FormatInt(timeNow, 10)

	urlStr := "http://tl.cyg.changyou.com/transaction/captcha-image"

	bodyMap := utils.BodyMap{}

	bodyMap.Set("goods_serial_num", "202406182116032545")

	bodyMap.Set("t", timeStr)

	_, body, err, i2 := utils.HttpRequest(context.Background(), "GET", urlStr+"?"+bodyMap.EncodeURLParams(), nil, nil)
	if err != nil {
		fmt.Printf("HttpRequest => %v\n", err)
		return
	}

	if i2 != 200 {
		fmt.Printf("请求失败，状态码：%d\n", i2)
		return
	}

	//保存图片
	path := "build/captcha.png"

	err = os.WriteFile(path, body, 0666)

	if err != nil {
		fmt.Printf("保存图片失败：%v\n", err)
		return
	}

	accessKeyId, accessKeySecret := "LTAI5tRpQYu9tz9gT6hVLPMT", "4ifEMUWxBe9JxOI36MmoefS1ssQkqG"

	client, err := ocr.NewAliOcrClient(accessKeyId, accessKeySecret)

	if err != nil {
		fmt.Printf("NewAliOcrClient => %v\n", err)
		return
	}

	verCode, err := client.OcrVerificationCode(path)

	if err != nil {
		fmt.Printf("OcrVerificationCode => %v\n", err)
		return
	}

	fmt.Printf("Verification code: %s\n", verCode)
}
