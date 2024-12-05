package mgr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lgrisa/lib/utils"
	"github.com/lgrisa/lib/utils/compress"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

func (d *excel_zip) GenerateSqlite() error {
	return d.doGenerate(true, false, true)
}

func (d *excel_zip) GenerateCs() error {
	return d.doGenerate(false, true, false)
}

func (d *excel_zip) doGenerate(shouldCache, createCS, createDB bool) error {
	fileMap, err := compress.UnpackZipData(d.data)
	if err != nil {
		return errors.Wrapf(err, "unpackData")
	}

	errorValue := atomic.Value{}
	wg := sync.WaitGroup{}

	for filename, fileBytes := range fileMap {
		if strings.Contains(filename, "__MACOSX") {
			continue
		}

		basename := filepath.Base(filename)

		if strings.HasPrefix(basename, "~$") {
			continue
		}

		if !strings.HasSuffix(basename, ".xlsx") {
			continue
		}

		dataMd5 := utils.Md5String(fileBytes)

		if shouldCache {
			// 从缓存中看下这个文件生成过没有，如果有，直接从缓存加载
			// {root}/{version}/{basename}/{md5}.json
			cachePath := fmt.Sprintf("%v/%v/%v/%v.json", d.root, d.version, basename, dataMd5)
			cacheBytes, err := os.ReadFile(cachePath)
			if err == nil {
				// 存在缓存
				file := &ExcelFile{}
				if err := json.Unmarshal(cacheBytes, file); err != nil {
					return errors.Wrapf(err, "jsonUnmarshal cachePath(%v) fail", cachePath)
				}

				d.fileMap[basename] = file
				continue
			}

			if !os.IsNotExist(err) {
				return errors.Wrapf(err, "read cachePath(%v) fail", cachePath)
			}
		}

		file := &ExcelFile{}
		file.Name = basename
		file.data = fileBytes
		file.DataMd5 = dataMd5
		file.DataSize = len(fileBytes)
		file.SheetMap = make(map[string]*ExcelSheet)

		d.fileMap[basename] = file

		wg.Add(1)
		go func(filename string, fileBytes []byte) {
			defer wg.Done()
			xlsxFile, err := xlsx.OpenBinary(fileBytes)
			if err != nil {
				err = errors.Wrapf(err, "xlsx.OpenBinary(%v) fail", filename)
				errorValue.Store(err)
				fmt.Println(err)
				return
			}
			file.xlsxFile = xlsxFile

			listSheet := xlsxFile.Sheet["list"]
			if listSheet == nil {
				err = errors.Errorf("表格数据中没有找到list页签 表名：%s ", filename)
				errorValue.Store(err)
				fmt.Println(err)
				return
			}

			for i := 0; i < len(listSheet.Rows); i++ {
				sheetRow := listSheet.Rows[i]

				if len(sheetRow.Cells) < 1 {
					continue
				}
				sheetName := strings.TrimSpace(sheetRow.Cells[0].String())

				curSheet := xlsxFile.Sheet[sheetName]
				if curSheet == nil {
					fmt.Println("找不到sheet 表名：[", filename, "] sheetName  [", sheetName, "]跳过")
					continue
				}

				_, fieldArray, err := newExcelSheetField(curSheet, d.idGen, filename, sheetName)
				if err != nil {
					err = errors.Wrapf(err, "表格数据解析失败 表名：%s ", filename)
					errorValue.Store(err)
					fmt.Println(err)
					return
				}

				es := &ExcelSheet{}
				es.Name = sheetName
				es.XlsxName = basename
				es.XlsxNameNoExt = strings.TrimSuffix(basename, ".xlsx")
				es.FullName = fmt.Sprintf("%v:%v", filename, sheetName)
				es.sheet = curSheet
				//es.fieldMap = fieldMap
				es.FieldArray = fieldArray

				es.ProtoMessageName = GetMessageName(es.XlsxNameNoExt, sheetName)
				es.protoBytes = es.newProtoBytes()
				es.ProtoMd5 = utils.Md5String(es.protoBytes)

				// {root}/{version}/{basename}/{sheet}_{proto_md5}.proto
				es.ProtoFilename = fmt.Sprintf("%v/%v/%v_%v.proto", d.version, file.Name, sheetName, es.ProtoMd5)
				es.rootProtoFilename = fmt.Sprintf("%v/%v", d.root, es.ProtoFilename)

				fileMd5, _ := utils.ReadFileMd5(es.rootProtoFilename)
				if fileMd5 != es.ProtoMd5 {
					// 如果proto文件不存在，就写到本地
					if err = os.WriteFile(es.rootProtoFilename, es.protoBytes, os.ModePerm); err != nil {
						errorValue.Store(errors.Wrapf(err, "写入Proto文件失败 表名：%s ", filename))
						return
					}
				}

				// 生成CS文件的路径（CS在请求的时候生成，这里不生成）
				es.rootCsDir = d.rootCsDir
				es.rootCsDirProto = fmt.Sprintf("%v/%v.proto", es.rootCsDir, es.ProtoMessageName)
				es.rootCsFilename = fmt.Sprintf("%v/%v.cs", es.rootCsDir, strings.Title(es.ProtoMessageName))
				if createCS {
					if err = es.CompileCs(); err != nil {
						err = errors.Wrapf(err, "编译Proto2Cs代码失败 表名：%s ", filename)
						errorValue.Store(err)
						fmt.Println(err)
						return
					}
				}

				// 生成Sqlite文件的路径（db文件的生成，调方法处理，这里是db和cs都会调用）
				// {root}/{version}/{basename}/sheet_{xlsx_md5}.db
				es.SqliteFilename = fmt.Sprintf("%v/%v_%v_%v.db", d.version, strings.TrimSuffix(file.Name, ".xlsx"), sheetName, file.DataMd5)
				es.rootSqliteFilename = fmt.Sprintf("%v/%v", d.root, es.SqliteFilename)

				// 如果是生成db请求，这里直接把db生成出来
				if createDB {
					if err := es.generateSqliteDB(d.drivers, d.storage); err != nil {
						err = errors.Wrapf(err, "生成SqliteDB失败 表名：%s ", filename)
						errorValue.Store(err)
						fmt.Println(err)
						return
					}
				}

				file.SheetMap[sheetName] = es
			}

			if shouldCache {
				data, err := json.Marshal(file)
				if err != nil {
					fmt.Printf("缓存file.json, json.Marshal(%v) fail, err: %+v \n", filename, err)
					return
				}

				cachePath := fmt.Sprintf("%v/%v/%v/%v.json", d.root, d.version, basename, dataMd5)
				if err = os.WriteFile(cachePath, data, os.ModePerm); err != nil {
					fmt.Printf("缓存file.json, writeFile(%v) fail, err: %+v \n", filename, err)
					return
				}
			}

		}(filename, fileBytes)
	}

	wg.Wait()

	if errRef := errorValue.Load(); errRef != nil {
		return errRef.(error)
	}

	if err := d.writeVersionFile(); err != nil {
		return errors.Wrapf(err, "写入版本文件失败")
	}

	return nil
}

func (d *excel_zip) packCsBody() ([]byte, error) {

	type protoData struct {
		ProtoFilename string
		ProtoBytes    []byte
	}

	var protoBuf []protoData

	fileBytes := make(map[string][]byte)
	for _, f := range d.fileMap {
		for _, s := range f.SheetMap {

			filename := filepath.Base(s.rootCsFilename)
			data, err := os.ReadFile(s.rootCsFilename)
			if err != nil {
				return nil, errors.Wrapf(err, "读取cs文件失败, %v", s.rootCsFilename)
			}
			fileBytes[filename] = data

			//protoFilename := filepath.Base(s.rootCsDirProto)
			//protoData, err := os.ReadFile(s.rootCsDirProto)
			//if err != nil {
			//	return nil, errors.Wrapf(err, "读取proto文件失败, %v", s.rootCsDirProto)
			//}
			//fileBytes[protoFilename] = protoData

			protoFilenameData, err := os.ReadFile(s.rootProtoFilename)
			if err != nil {
				return nil, errors.Wrapf(err, "读取proto文件失败, %v", s.rootProtoFilename)
			}

			//处理掉开头
			protoFilenameData = bytes.TrimPrefix(protoFilenameData, []byte(ProtoPrefix))

			protoBuf = append(protoBuf, protoData{
				ProtoFilename: filepath.Base(s.rootProtoFilename),
				ProtoBytes:    protoFilenameData,
			})
		}
	}

	//整理排序
	sort.Slice(protoBuf, func(i, j int) bool {
		return protoBuf[i].ProtoFilename < protoBuf[j].ProtoFilename
	})

	var protoBufBytes bytes.Buffer

	protoBufBytes.WriteString(ProtoPrefix)
	protoBufBytes.WriteString("\n")
	for _, v := range protoBuf {
		protoBufBytes.Write(v.ProtoBytes)
	}

	fileBytes["conf.proto"] = protoBufBytes.Bytes()

	return compress.PackZipData(fileBytes)
}
