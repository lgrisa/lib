package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func GetFilePrefix(fileName string) string {
	fileNameAll := path.Base(fileName)
	fileSuffix := path.Ext(fileNameAll)
	filePrefix := fileNameAll[0 : len(fileNameAll)-len(fileSuffix)]

	return filePrefix
}

func ReplaceFileContentField(filePath string, replaceField string, insteadFile string) error {

	in, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		_ = in.Close()
	}(in)

	out, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}

	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	br := bufio.NewReader(in)
	index := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		newLine := strings.Replace(string(line), replaceField, insteadFile, -1)
		_, err = out.WriteString(newLine + "\n")
		if err != nil {
			fmt.Println("write to file fail:", err)
			os.Exit(-1)
		}
		index++
	}

	return nil
}
