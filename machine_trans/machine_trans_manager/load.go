package machine_trans_manager

import (
	"fmt"
	"github.com/lgrisa/lib/machine_trans/machine_trans_engine"
	"github.com/lgrisa/lib/utils/logutil"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
)

func (m *Manager) initConfig() error {

	file, err := xlsx.OpenFile(m.configPath)
	if err != nil {
		return errors.Errorf("initConfig: open file: %s error: %v", m.configPath, err)
	}

	configSheet, ok := file.Sheet["LanguageConfig"]

	if !ok {
		return errors.Errorf("initConfig: %v not found sheet LanguageConfig", m.configPath)
	}

	// 读取单元格数据
	for i := 4; i < len(configSheet.Rows); i++ {

		currentRow := configSheet.Rows[i]

		if len(currentRow.Cells) < 4 {
			logutil.LogErrorF(fmt.Sprintf("initConfig: row %d cells length < 3", i))
			continue
		}

		fileN := currentRow.Cells[0].String()

		isOpen := currentRow.Cells[3].String()

		if isOpen != "1" {
			logutil.LogTraceF(fmt.Sprintf("initConfig: file %s is not open", fileN))
			continue
		}

		dirPath := m.excelPath + "/" + fileN
		// 检查目录是否存在
		_, errExist := os.Stat(dirPath)

		if errExist != nil {
			if os.IsNotExist(errExist) {
				return errors.Errorf("initConfig: file %s not exist", dirPath)
			} else {
				return errors.Errorf("initConfig: file %s stat error: %v", dirPath, errExist)
			}
		}

		languageType, errGet := machine_trans_engine.GetLanguageType(fileN)

		if errGet != nil {
			return errors.Errorf("trans.GetLanguageType => %v", errGet)
		}

		m.TransConfigMap[languageType] = &TranslatedLanguageInfo{
			languageType: languageType,
			path:         dirPath,
		}

		logutil.LogInfoF(fmt.Sprintf("加载翻译配置表: 对应语言 %s 需要进行翻译", fileN))
	}

	logutil.LogInfoF(fmt.Sprintf("加载翻译配置表:%v 成功", m.configPath))

	return nil
}

func (m *Manager) LoadOriginExcel() error {
	files, err := os.ReadDir(m.excelPath)
	if err != nil {
		return errors.Errorf("LoadOriginExcel: read dir: %s error: %v", m.excelPath, err)
	}

	var waitGroup sync.WaitGroup
	errAtomic := atomic.Value{}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		filename := file.Name()

		if !strings.HasPrefix(filename, "Local") {
			continue
		}

		m.OriginalExcelMap[filename] = make(map[string]*OriginalExcelInfo)

		currentFile := file

		waitGroup.Add(1)
		go func(file os.DirEntry) {

			defer func() {
				waitGroup.Done()

				if r := recover(); r != nil {
					// 打印堆栈

					logutil.LogErrorF(fmt.Sprintf("加载原始表: %v error: %v stack: %v", file.Name(), r, debug.Stack()))
					errAtomic.Store(r)
				}

			}()

			// 读取文件
			if errLoad := m.loadExcel(file); errLoad != nil {
				errAtomic.Store(errLoad)
				return
			}

		}(currentFile)
	}

	waitGroup.Wait()

	if errWait := errAtomic.Load(); errWait != nil {
		return errors.Errorf("LoadOriginExcel: error: %v", errWait)
	}

	logutil.LogInfoF(fmt.Sprintf("加载原始表：%v 成功", m.excelPath))

	return nil
}

func (m *Manager) loadExcel(file os.DirEntry) error {
	// 解析文件

	filename := m.excelPath + "/" + file.Name()

	f, err := xlsx.OpenFile(filename)

	if err != nil {
		return errors.Errorf("LoadOriginExcel: open file: %s error: %v", filename, err)
	}

	//找到对应Data sheet
	sheet, ok := f.Sheet["Localization"]

	if !ok {
		return errors.Errorf("LoadOriginExcel: %s not found sheet Localization", filename)
	}

	curMap := m.OriginalExcelMap[file.Name()]

	// 读取单元格数据
	for i := 3; i < len(sheet.Rows); i++ {

		currentRow := sheet.Rows[i]

		if len(currentRow.Cells) == 0 {
			continue
		}

		if len(currentRow.Cells) < 2 {
			logutil.LogTraceF(fmt.Sprintf("LoadOriginExcel: file %s row %d cells length < 2", filename, i))
			continue
		}

		keyID := currentRow.Cells[0].String()
		valueCn := currentRow.Cells[1].String()

		if keyID == "" {
			logutil.LogTraceF(fmt.Sprintf("LoadOriginExcel: file %s row %d keyID is empty", filename, i))
			continue
		}

		// 以//开头的行不处理
		if strings.HasPrefix(keyID, "//") {
			continue
		}

		if valueCn == "" {
			logutil.LogTraceF(fmt.Sprintf("解析原始文件: file %s row %d 中文不存在跳过", filename, i))
			continue
		}

		// 保存到map
		curMap[keyID] = &OriginalExcelInfo{
			keyID:   keyID,
			valueCn: valueCn,
		}
	}

	logutil.LogInfoF(fmt.Sprintf("加载原始中文表: 文件名称 %s 加载翻译数量: %d", filename, len(curMap)))

	return nil
}
