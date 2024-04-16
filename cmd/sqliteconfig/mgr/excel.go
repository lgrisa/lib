package mgr

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"sort"
	"strings"
	"time"
)

func newExcelZip(data []byte, md5 string, root, version string, idGen *MessageIdGen, drivers *Drivers, storage *Storage) *excel_zip {

	rootCsDir := fmt.Sprintf("%v/cs/%v", root, time.Now().UnixNano())

	return &excel_zip{
		data:      data,
		dataMd5:   md5,
		root:      root,
		version:   version,
		rootCsDir: rootCsDir,
		idGen:     idGen,
		drivers:   drivers,
		storage:   storage,
		fileMap:   make(map[string]*ExcelFile),
	}
}

type excel_zip struct {
	data    []byte
	dataMd5 string

	root    string
	version string

	// {root}/cs/{unixtimestamp}
	rootCsDir string

	// {root}/proto_id.yaml
	idGen   *MessageIdGen
	drivers *Drivers

	storage *Storage

	fileMap map[string]*ExcelFile
}

type ExcelFile struct {
	Name     string
	data     []byte
	DataMd5  string
	DataSize int

	xlsxFile *xlsx.File

	SheetMap map[string]*ExcelSheet
}

type ExcelSheet struct {
	// data
	Name string

	// hero.xlsx
	XlsxName string

	// hero
	XlsxNameNoExt string

	// hero.xlsx:data
	FullName string

	sheet *xlsx.Sheet

	//fieldMap map[string]*ExcelSheetField

	FieldArray []*ExcelSheetField

	// 生成的proto消息
	ProtoMessageName string

	protoBytes []byte

	ProtoMd5 string

	// 因为sheet并没有单独的md5，只能使用excel的md5
	// 只要excel有变化，内部所有sheet都需要重新生成

	// 生成Proto文件的路径
	// {root}/{version}/{basename}/sheet_{proto_md5}.proto
	rootProtoFilename string
	ProtoFilename     string

	// 生成CS文件的路径
	// {root}/cs/{unixtimestamp}
	rootCsDir string
	// {rootCsDir}/{ProtoMessageName}.proto
	rootCsDirProto string
	// {rootCsDir}/{ProtoMessageName}.cs
	rootCsFilename string

	// 生成Sqlite文件的路径
	// {root}/{version}/{basename}/sheet_{xlsx_md5}.db
	rootSqliteFilename string
	SqliteFilename     string
}

type ExcelSheetField struct {
	Name string

	Type string

	ProtoType string

	IsRepeated bool

	ProtoId int
}

const starReadLine = 5

func newExcelSheetField(curSheet *xlsx.Sheet, idGen *MessageIdGen, fileName, sheetName string) (map[string]*ExcelSheetField, []*ExcelSheetField, error) {
	if len(curSheet.Rows) < starReadLine {
		return nil, nil, errors.Errorf("表格行数不足跳过 表名：%v 页签名称：%v 行数：%d ，最低行数要求：%d \n", fileName, sheetName, len(curSheet.Rows), starReadLine)
	}

	titleRow := curSheet.Rows[1]
	typeRow := curSheet.Rows[2]

	fieldMap := make(map[string]*ExcelSheetField)
	var fieldArray []*ExcelSheetField
	for j := 0; j < len(titleRow.Cells); j++ {

		title := titleRow.Cells[j].String()
		if title == "" {
			continue
		}

		if f := fieldMap[title]; f != nil {
			// 字段已经存在，那么属于是多个字段来表示数组
			f.IsRepeated = true
			continue
		}

		titleType := strings.ToLower(getCellString(typeRow, j))
		if titleType == "" {
			fmt.Printf("表格列类型为空跳过 表名：%v 页签名称：%v 列数：%d \n", fileName, sheetName, j)
			continue
		}

		protoType, isRepeated := getProtoType(titleType) //不会出现NULL

		f := &ExcelSheetField{}
		f.Name = title
		f.Type = titleType
		f.IsRepeated = isRepeated
		f.ProtoType = protoType

		fieldMap[title] = f
		fieldArray = append(fieldArray, f)
	}

	for _, f := range fieldArray {
		fieldType := strings.TrimSpace(f.ProtoType)
		if f.IsRepeated {
			fieldType = fieldType + "_array"
		}

		f.ProtoId = idGen.MsgProtoFieldId(fileName, sheetName, f.Name, fieldType)
	}

	sort.Slice(fieldArray, func(i, j int) bool {
		return fieldArray[i].ProtoId < fieldArray[j].ProtoId
	})

	return fieldMap, fieldArray, nil
}

func getCellString(row *xlsx.Row, i int) string {
	if i < len(row.Cells) {
		return strings.TrimSpace(row.Cells[i].String())
	}
	return ""
}
