package mgr

import (
	"fmt"
	"github.com/lgrisa/lib/utils"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (m *Manager) handleGenSqlite(w http.ResponseWriter, r *http.Request) {
	// 从header中获取md5值
	fileMd5 := r.Header.Get(ReqHeaderFileMd5)
	fileSize := r.Header.Get(ReqHeaderFileSize)

	if fileMd5 == "" || fileSize == "" {
		writeErrMsg(w, "md5 or size is empty")
		return
	}

	var fileSizeInt64 int64
	val, err := strconv.ParseInt(fileSize, 10, 64)
	if err != nil || val <= 0 {
		writeErrMsg(w, "invalid size")
		return
	}
	fileSizeInt64 = val

	utils.LogTraceF("handleGenSqlite fileMd5: (%v), fileSize: (%v)", fileMd5, fileSize)

	confJson, err := loadConfigJson(m.root, packVersion, fileMd5, int(fileSizeInt64))
	if err != nil {
		if !os.IsNotExist(err) {
			writeErrMsg(w, fmt.Sprintf("loadConfigJson fail: %v", err))
			return
		}
	} else {
		writeConfigJson(w, confJson, "使用已存在对应版本:"+fileMd5)
		return
	}

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

	if err = m.funcIdMap(func(idMap *MessageIdGen, drivers *Drivers) bool {

		excelZip := newExcelZip(data, dataMd5, m.root, packVersion, idMap, drivers, m.storage)
		if err = excelZip.GenerateSqlite(); err != nil {
			errMsg := fmt.Sprintf("生成Sqlite文件失败, %v", err)
			writeErrMsg(w, errMsg)
			return false
		}

		confJson, err = loadConfigJson(m.root, packVersion, fileMd5, int(fileSizeInt64))
		if err != nil {
			writeErrMsg(w, "loadConfigJson fail: "+err.Error())
			return false
		} else {
			writeConfigJson(w, confJson, "使用生成的版本:"+fileMd5)
			return true
		}
	}); err != nil {
		writeErrMsg(w, fmt.Sprintf("funcIdMap fail, %v", err))
		return
	}
}
