package machine_trans_manager

import (
	"fmt"
	"github.com/lgrisa/lib/machine_trans/machine_trans_engine"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
)

func (m *Manager) translate() error {

	for _, languageInfo := range m.TransConfigMap {

		files, err := os.ReadDir(m.excelPath)

		if err != nil {
			return errors.Errorf("translate: read dir: %s error: %v", m.excelPath, err)
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

			currentFile := file

			waitGroup.Add(1)
			go func(file os.DirEntry) {

				defer func() {
					waitGroup.Done()

					if r := recover(); r != nil {
						// 打印堆栈

						utils.LogErrorF(fmt.Sprintf("translate:(%v) error: %v stack: %v", file.Name(), r, string(debug.Stack())))
						errAtomic.Store(r)
					}

				}()

				// 读取文件
				if errLoad := m.translateLanguageExcel(file, languageInfo); errLoad != nil {
					errAtomic.Store(errLoad)
					return
				}

			}(currentFile)
		}

		waitGroup.Wait()

		if errWait := errAtomic.Load(); errWait != nil {
			return errors.Errorf("translate: error: %v", errWait)
		}
	}

	utils.LogInfoF("翻译结束")

	return nil
}

var transValueCellIndex = 2

func (m *Manager) translateLanguageExcel(file os.DirEntry, languageInfo *TranslatedLanguageInfo) error {
	// 解析文件
	filename := languageInfo.path + "/" + file.Name()

	utils.LogInfoF(fmt.Sprintf("开始处理表格: file: %v", filename))

	f, err := xlsx.OpenFile(filename)

	if err != nil {
		return errors.Errorf("translateLanguageExcel: file: %s error: %v", filename, err)
	}

	//找到对应Data sheet
	sheet, ok := f.Sheet["Localization"]

	if !ok {
		return errors.Errorf("translateLanguageExcel: file %s not found sheet Localization", filename)
	}

	curOriginMap := m.OriginalExcelMap[file.Name()]

	if curOriginMap == nil {
		return errors.Errorf("translateLanguageExcel: file %s not found original map", filename)
	}

	curExcelMap := make(map[string]*OriginalExcelInfo)

	isNeedSave := false
	directTransCount := 0

	for i := 3; i < len(sheet.Rows); i++ {

		currentRow := sheet.Rows[i]

		if len(currentRow.Cells) == 0 {
			continue
		}

		if len(currentRow.Cells) < 2 {
			utils.LogTraceF(fmt.Sprintf("translateLanguageExcel: file %s row %d cells length < 2", filename, i))
			continue
		}

		keyID := currentRow.Cells[0].String()
		valueLocal := currentRow.Cells[1].String()

		if keyID == "" {
			utils.LogTraceF(fmt.Sprintf("translateLanguageExcel: file %s row %d keyID is empty", filename, i))
			continue
		}

		if valueLocal == "" {
			utils.LogTraceF(fmt.Sprintf("translateLanguageExcel: file %s row %d keyID: %s 中文不存在,跳过", filename, i, keyID))
			continue
		}

		transValue := ""

		if len(currentRow.Cells) > transValueCellIndex {
			transValue = currentRow.Cells[transValueCellIndex].String()
		}

		if transValue == "" {

			// 直接翻译
			directTransCount++

			transValue, err = m.engineTranslateFor(valueLocal, machine_trans_engine.LanguageTypeZh, languageInfo.languageType)

			if err != nil {
				return errors.Errorf("未翻译 转义接口报错: 文件名：%s keyID: %s error: %v", filename, keyID, err)
			}

			// 保存到文件

			transValue = getTransValue(transValue)

			if len(currentRow.Cells) <= transValueCellIndex {
				currentRow.AddCell().SetValue(transValue)
			} else {
				currentRow.Cells[transValueCellIndex].SetValue(transValue)
			}

			isNeedSave = true
		}

		// 保存到map
		curExcelMap[keyID] = &OriginalExcelInfo{
			keyID:      keyID,
			valueLocal: valueLocal,
		}
	}

	utils.LogInfoF(fmt.Sprintf("对应语言：%v file %s 原始中文表大小: %d, 当前对于语言表大小: %d 直接翻译数量: %d", languageInfo.languageType, filename, len(curOriginMap), len(curExcelMap), directTransCount))

	newAddMap := make(map[string]*OriginalExcelInfo)
	// 处理新增
	for key, value := range curOriginMap {
		if _, isExist := curExcelMap[key]; !isExist {
			newAddMap[key] = value
		}
	}

	if len(newAddMap) > 0 {
		isNeedSave = true

		for key, value := range newAddMap {
			row := sheet.AddRow()
			row.AddCell().SetValue(key)
			row.AddCell().SetValue(value.valueCn)

			// 直接翻译
			transValue, errTrans := m.engineTranslateFor(value.valueCn, "zh", languageInfo.languageType)

			if errTrans != nil {
				return errors.Errorf("新增 转义接口报错: 文件名：%s keyID: %s 中文: %s error: %v", filename, key, value.valueCn, errTrans)
			}

			row.AddCell().SetValue(getTransValue(transValue))
		}
	}

	utils.LogInfoF(fmt.Sprintf("对应语言：%v file %s 新增翻译数量: %d", languageInfo.languageType, filename, len(newAddMap)))

	// 保存到文件
	if isNeedSave {
		errSave := f.Save(filename)

		if errSave != nil {
			return errors.Errorf("translateLanguageExcel: file %s save error: %v", filename, errSave)
		}
	}

	return nil
}

func (m *Manager) engineTranslateFor(text string, fromLanguage, toLanguage machine_trans_engine.LanguageType) (string, error) {

	curEngine, isExist := m.TransClient[HuaWei]

	if isExist {
		if resp, err := curEngine.TranslateFor(text, fromLanguage, toLanguage); err != nil {
			if !errors.Is(err, machine_trans_engine.ErrLanguageTypeNotSupported) {
				return "", err
			}

			// 如果不支持的语言类型，尝试使用其他引擎
		} else {
			return resp, nil
		}
	}

	firstEngine, isExist := m.TransClient[VolcEngine]

	if isExist {
		return firstEngine.TranslateFor(text, fromLanguage, toLanguage)
	}

	return "", errors.New("无可用翻译引擎")
}

func getTransValue(value string) string {
	return "@@" + value
}
