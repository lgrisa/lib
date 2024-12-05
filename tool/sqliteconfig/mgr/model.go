package mgr

import (
	"github.com/lgrisa/lib/tool/sqliteconfig/s3"
	"sync"
)

type Manager struct {
	port int
	root string

	idMapMux  sync.Mutex
	idMapPath string

	drivers *Drivers
	storage *s3.Storage
}

const packVersion = "1"

const (
	ReqHeaderFileMd5  = "file_md5"
	ReqHeaderFileSize = "file_size"
)

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

type ConfigCsJson struct {
	FileMd5 string `json:"file_md5"`

	FileSize int64 `json:"file_size"`

	CsBody []byte `json:"cs_body"`
}
