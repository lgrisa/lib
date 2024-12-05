package main

import (
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
)

// 生成路径

var GeneratePath = flag.String("GeneratePath", "", "生成路径")

// 针对的url

var url = flag.String("url", "", "url")

func main() {
	flag.Parse()

	if *GeneratePath == "" {
		fmt.Println("please input path")
		return
	}

	if *url == "" {
		fmt.Println("please input url")
		return
	}

	if err := qrcode.WriteFile(*url, qrcode.Medium, 256, *GeneratePath); err != nil {
		fmt.Printf("生成二维码失败: %s, err: %s", *GeneratePath, err)
		return
	}
}
