package utils

import "path"

func GetFilePrefix(fileName string) string {
	fileNameAll := path.Base(fileName)
	fileSuffix := path.Ext(fileNameAll)
	filePrefix := fileNameAll[0 : len(fileNameAll)-len(fileSuffix)]

	return filePrefix
}
