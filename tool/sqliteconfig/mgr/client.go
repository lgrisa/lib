package mgr

import (
	"fmt"
	"github.com/lgrisa/lib/utils"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (m *Manager) handleGenCs(w http.ResponseWriter, r *http.Request) {
	// 从header中获取md5值
	fileMd5 := r.Header.Get(ReqHeaderFileMd5)
	fileSize := r.Header.Get(reqHeaderFileSize)

	if fileMd5 == "" || fileSize == "" {
		fmt.Println("md5 or size is empty")
		writeErrMsg(w, "md5 or size is empty")
		return
	}

	var fileSizeInt64 int64
	val, err := strconv.ParseInt(fileSize, 10, 64)
	if err != nil || val <= 0 {
		writeErrMsg(w, fmt.Sprintf("invalid file size, %v, %v", fileSize, val))
		return
	}
	fileSizeInt64 = val

	fmt.Println("handleGenCs", fileMd5, fileSize)

	//重新生成
	data, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrMsg(w, "io.ReadAll(r.Body) fail: "+err.Error())
		return
	}

	dataMd5 := utils.Md5String(data)
	dataSize := int64(len(data))

	if fileMd5 != dataMd5 || fileSizeInt64 != dataSize {
		writeErrMsg(w, fmt.Sprintf("header和body的数据不一致, headerMd5: %v, bodyMd5: %v, headerSize: %v, bodySize: %v", fileMd5, dataMd5, fileSizeInt64, dataSize))
		return
	}

	if err := m.funcIdMap(func(idMap *MessageIdGen, drivers *Drivers) bool {

		excelZip := newExcelZip(data, dataMd5, m.root, packVersion, idMap, drivers, m.storage)
		defer func() {
			_ = os.RemoveAll(excelZip.rootCsDir)
		}()

		if err := excelZip.GenerateCs(); err != nil {
			errMsg := fmt.Sprintf("生成Cs文件失败, %v", err)
			fmt.Println(errMsg)
			writeErrMsg(w, errMsg)
			return false
		}

		csBody, err := excelZip.packCsBody()
		if err != nil {
			errMsg := fmt.Sprintf("打包csBody文件失败, %v", err)
			fmt.Println(errMsg)
			writeErrMsg(w, errMsg)
			return false
		}

		confJson := &ConfigCsJson{}
		confJson.FileMd5 = dataMd5
		confJson.FileSize = dataSize
		confJson.CsBody = csBody

		writeCsJson(w, confJson, "使用生成的版本:"+fileMd5)
		return true
	}); err != nil {
		errMsg := fmt.Sprintf("funcIdMap fail, %v", err)
		writeErrMsg(w, errMsg)
		return
	}
}
