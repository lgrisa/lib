package main

import (
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
)

var path = flag.String("path", "", "path")

var url = flag.String("url", "", "url")

func main() {
	flag.Parse()

	if *path == "" {
		fmt.Println("please input path")
		return
	}

	if *url == "" {
		fmt.Println("please input url")
		return
	}

	if err := qrcode.WriteFile(*url, qrcode.Medium, 256, *path); err != nil {
		fmt.Printf("生成二维码失败: %s, err: %s", *path, err)
		return
	}
}
