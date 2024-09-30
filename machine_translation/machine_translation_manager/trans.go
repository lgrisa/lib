package machine_translation_manager

import (
	"fmt"
	"github.com/lgrisa/lib/machine_translation/machine_translation_engine"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func NewClient(configPath string, excelPath string) (*Manager, error) {

	m := &Manager{
		TransConfigMap:   make(map[machine_translation_engine.LanguageType]*TranslatedLanguageInfo),
		OriginalExcelMap: make(map[string]map[string]*OriginalExcelInfo),
		TransClient:      make(map[EngineType]TranslatedClient),
		excelPath:        excelPath,
		configPath:       configPath,
	}

	return m, nil
}

func (m *Manager) RegisterClient(engineType EngineType, client TranslatedClient) error {

	if client == nil {
		return errors.Errorf("RegisterClient:%v client is nil", engineType)
	}

	m.TransClient[engineType] = client

	return nil
}

func (m *Manager) Start() error {

	if errInitConfig := m.initConfig(); errInitConfig != nil {
		return errors.Errorf("m.initConfig => %v", errInitConfig)
	}

	if errLoadOriginExcel := m.LoadOriginExcel(); errLoadOriginExcel != nil {
		return errors.Errorf("m.LoadOriginExcel => %v", errLoadOriginExcel)
	}

	if errTranslate := m.startTranslate(); errTranslate != nil {
		return errors.Errorf("m.startTranslate => %v", errTranslate)
	}

	return nil
}

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
			utils.LogErrorF(fmt.Sprintf("initConfig: row %d cells length < 3", i))
			continue
		}

		fileN := currentRow.Cells[0].String()

		isOpen := currentRow.Cells[3].String()

		if isOpen != "1" {
			utils.LogTraceF(fmt.Sprintf("initConfig: file %s is not open", fileN))
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

		languageType, errGet := machine_translation_engine.GetLanguageType(fileN)

		if errGet != nil {
			return errors.Errorf("trans.GetLanguageType => %v", errGet)
		}

		m.TransConfigMap[languageType] = &TranslatedLanguageInfo{
			languageType: languageType,
			path:         dirPath,
		}

		utils.LogInfoF(fmt.Sprintf("加载翻译配置表: 对应语言 %s 需要进行翻译", fileN))
	}

	utils.LogInfoF(fmt.Sprintf("加载翻译配置表:%v 成功", m.configPath))

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
					utils.LogErrorF(fmt.Sprintf("Error reading file: %v", r))
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

	utils.LogInfoF(fmt.Sprintf("加载原始表：%v 成功", m.excelPath))

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
			utils.LogTraceF(fmt.Sprintf("LoadOriginExcel: file %s row %d cells length < 2", filename, i))
			continue
		}

		keyID := currentRow.Cells[0].String()
		valueCn := currentRow.Cells[1].String()

		if keyID == "" {
			utils.LogTraceF(fmt.Sprintf("LoadOriginExcel: file %s row %d keyID is empty", filename, i))
			continue
		}

		// 以//开头的行不处理
		if strings.HasPrefix(keyID, "//") {
			continue
		}

		if valueCn == "" {
			utils.LogTraceF(fmt.Sprintf("解析原始文件: file %s row %d 中文不存在跳过", filename, i))
			continue
		}

		// 保存到map
		curMap[keyID] = &OriginalExcelInfo{
			keyID:   keyID,
			valueCn: valueCn,
		}
	}

	utils.LogInfoF(fmt.Sprintf("加载原始中文表: 文件名称 %s 加载翻译数量: %d", filename, len(curMap)))

	return nil
}

func (m *Manager) startTranslate() error {

	for _, languageInfo := range m.TransConfigMap {

		files, err := os.ReadDir(m.excelPath)

		if err != nil {
			return errors.Errorf("startTranslate: read dir: %s error: %v", m.excelPath, err)
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
						utils.LogErrorF(fmt.Sprintf("Error reading file: %v", r))
						errAtomic.Store(r)
					}

				}()

				// 读取文件
				if errLoad := m.startTranslateLanguageExcel(file, languageInfo); errLoad != nil {
					errAtomic.Store(errLoad)
					return
				}

			}(currentFile)
		}

		waitGroup.Wait()

		if errWait := errAtomic.Load(); errWait != nil {
			return errors.Errorf("startTranslate: error: %v", errWait)
		}
	}

	utils.LogInfoF("翻译结束")

	return nil
}

var transValueCellIndex = 2

func (m *Manager) startTranslateLanguageExcel(file os.DirEntry, languageInfo *TranslatedLanguageInfo) error {
	// 解析文件

	filename := languageInfo.path + "/" + file.Name()

	f, err := xlsx.OpenFile(filename)

	if err != nil {
		return errors.Errorf("startTranslateLanguageExcel: file: %s error: %v", filename, err)
	}

	//找到对应Data sheet
	sheet, ok := f.Sheet["Localization"]

	if !ok {
		return errors.Errorf("startTranslateLanguageExcel: file %s not found sheet Localization", filename)
	}

	curOriginMap := m.OriginalExcelMap[file.Name()]

	if curOriginMap == nil {
		return errors.Errorf("startTranslateLanguageExcel: file %s not found original map", filename)
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
			utils.LogTraceF(fmt.Sprintf("startTranslateLanguageExcel: file %s row %d cells length < 2", filename, i))
			continue
		}

		keyID := currentRow.Cells[0].String()
		valueLocal := currentRow.Cells[1].String()

		if keyID == "" {
			utils.LogTraceF(fmt.Sprintf("startTranslateLanguageExcel: file %s row %d keyID is empty", filename, i))
			continue
		}

		if valueLocal == "" {
			utils.LogTraceF(fmt.Sprintf("startTranslateLanguageExcel: file %s row %d keyID: %s 中文不存在,跳过", filename, i, keyID))
			continue
		}

		transValue := ""

		if len(currentRow.Cells) > transValueCellIndex {
			transValue = currentRow.Cells[transValueCellIndex].String()
		}

		if transValue == "" {

			// 直接翻译
			directTransCount++

			transValue, err = m.TranslateFor(valueLocal, "zh", languageInfo.languageType)

			if err != nil {
				return errors.Errorf("未翻译 转义接口报错: 文件名：%s keyID: %s error: %v", filename, keyID, err)
			}

			// 保存到文件

			transValue = GetTransValue(transValue)

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
			transValue, errTrans := m.TranslateFor(value.valueCn, "zh", languageInfo.languageType)

			if errTrans != nil {
				return errors.Errorf("新增 转义接口报错: 文件名：%s keyID: %s 中文: %s error: %v", filename, key, value.valueCn, errTrans)
			}

			row.AddCell().SetValue(GetTransValue(transValue))
		}
	}

	utils.LogInfoF(fmt.Sprintf("对应语言：%v file %s 新增翻译数量: %d", languageInfo.languageType, filename, len(newAddMap)))

	// 保存到文件
	if isNeedSave {
		errSave := f.Save(filename)

		if errSave != nil {
			return errors.Errorf("startTranslateLanguageExcel: file %s save error: %v", filename, errSave)
		}
	}

	return nil
}

func (m *Manager) TranslateFor(text string, fromLanguage, toLanguage machine_translation_engine.LanguageType) (string, error) {

	firstEngine, isExist := m.TransClient[VolcEngine]

	if isExist {
		return firstEngine.TranslateFor(text, fromLanguage, toLanguage)
	}

	//循环所有引擎
	for _, engine := range m.TransClient {
		if res, err := engine.TranslateFor(text, fromLanguage, toLanguage); err == nil {
			return res, nil
		}
	}

	return "", errors.New("无可用翻译引擎")
}

func GetTransValue(value string) string {
	return "@@" + value
}
