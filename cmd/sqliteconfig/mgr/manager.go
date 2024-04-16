package mgr

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lgrisa/library/utils"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const packVersion = "1"

func NewManager(port int, root, idMapPath string, storage *Storage) *Manager {
	return &Manager{
		port:      port,
		root:      root,
		idMapPath: idMapPath,
		drivers:   NewDrivers(),
		storage:   storage,
	}
}

type Manager struct {
	port int

	root string

	idMapMux  sync.Mutex
	idMapPath string

	drivers *Drivers

	storage *Storage
}

func (m *Manager) Start() {

	http.HandleFunc("/server/", m.handleGenSqlite)
	http.HandleFunc("/client/cs/", m.handleGenCs)

	http.HandleFunc("/client/ts/", m.handleGenTypeScript)

	http.Handle("/client/sqlite/", http.StripPrefix("/client/sqlite/", http.FileServer(http.Dir(m.root))))

	fmt.Println("listen port:", m.port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server start fail", err)
		}
	}
	fmt.Println("server exit")
}

func (m *Manager) funcIdMap(f func(idMap *MessageIdGen, drivers *Drivers) bool) error {
	m.idMapMux.Lock()
	defer m.idMapMux.Unlock()

	// 加载上来
	idMap, err := loadGen(m.idMapPath)
	if err != nil {
		return errors.Wrapf(err, "加载idMap失败, %v", m.idMapPath)
	}

	// 使用
	if !f(idMap, m.drivers) {
		return nil
	}

	// 用完保存
	toSave, err := idMap.encode()
	if err != nil {
		return errors.Wrapf(err, "idMap.encode失败")
	}

	if err = os.WriteFile(m.idMapPath, toSave, os.ModePerm); err != nil {
		return errors.Wrapf(err, "写入idMap文件失败, %v", m.idMapPath)
	}

	return nil
}

func (m *Manager) handleGenSqlite(w http.ResponseWriter, r *http.Request) {
	// 从header中获取md5值
	fileMd5 := r.Header.Get("file_md5")
	fileSize := r.Header.Get("file_size")

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

	fmt.Println("handleGenSqlite", fileMd5, fileSize)

	confJson, err := loadConfigJson(m.root, packVersion, fileMd5, int(fileSizeInt64))
	if err != nil {
		if !os.IsNotExist(err) {
			writeErrMsg(w, "loadConfigJson fail: "+err.Error())
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

func (m *Manager) handleGenCs(w http.ResponseWriter, r *http.Request) {
	// 从header中获取md5值
	fileMd5 := r.Header.Get("file_md5")
	fileSize := r.Header.Get("file_size")

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

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	// code == 200时，返回的数据
	Data interface{} `json:"data"`
}

type ConfigJson struct {
	FileMd5 string `json:"file_md5"`

	FileSize int64 `json:"file_size"`

	VersionPath string `json:"version_path"`

	VersionMd5 string `json:"version_md5"`
}

func writeConfigJson(w http.ResponseWriter, data *ConfigJson, cmdOutPut string) {
	jsoniter.NewEncoder(w).Encode(&response{
		Code:    200,
		Data:    data,
		Message: cmdOutPut,
	})
}

func writeErrMsg(w http.ResponseWriter, s string) {
	_ = jsoniter.NewEncoder(w).Encode(&response{
		Code:    400,
		Message: s,
	})

	fmt.Println(s)
}

type ConfigCsJson struct {
	FileMd5 string `json:"file_md5"`

	FileSize int64 `json:"file_size"`

	CsBody []byte `json:"cs_body"`
}

func writeCsJson(w http.ResponseWriter, data *ConfigCsJson, cmdOutPut string) {
	jsoniter.NewEncoder(w).Encode(&response{
		Code:    200,
		Data:    data,
		Message: cmdOutPut,
	})
}
