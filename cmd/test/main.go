package main

import (
	"fmt"
	"path"
	"path/filepath"
)

func main() {
	path0, err := filepath.Abs(".")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(path0)

	pathDir := path.Dir(path0)

	fmt.Println(pathDir)

	pathBase := path.Join(path0, "conf")

	fmt.Println(pathBase)
}
