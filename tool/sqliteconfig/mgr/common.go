package mgr

import (
	"fmt"
	"github.com/lgrisa/lib/tool/sqliteconfig/s3"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func NewManager(port int, root, idMapPath string, storage *s3.Storage) *Manager {
	return &Manager{
		port:      port,
		root:      root,
		idMapPath: idMapPath,
		drivers:   NewDrivers(),
		storage:   storage,
	}
}

func (m *Manager) Register() {

	http.HandleFunc("/server/", m.handleGenSqlite)
	http.HandleFunc("/client/cs/", m.handleGenCs)
	http.HandleFunc("/client/ts/", m.handleGenTypeScript)
	http.Handle("/client/sqlite/", http.StripPrefix("/client/sqlite/", http.FileServer(http.Dir(m.root))))

	utils.LogTraceF("server start at %d", m.port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server start fail", err)
		}
	}

	utils.LogTraceF("server stop")
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
