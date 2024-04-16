package mgr

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewDrivers() *Drivers {
	return &Drivers{
		dataMap: make(map[string]*Driver),
	}
}

type Drivers struct {
	dataMap map[string]*Driver
}

func (d *Drivers) getOrCreate(name string) *Driver {
	if v := d.dataMap[name]; v != nil {
		return v
	}
	v := newDriver(name)
	d.dataMap[name] = v
	return v
}

func newDriver(name string) *Driver {

	sr := &Driver{}
	sr.name = name

	sql.Register(sr.name, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			sr.driverCons = append(sr.driverCons, conn)
			return nil
		},
	})

	return sr
}

type Driver struct {
	name       string
	driverCons []*sqlite3.SQLiteConn
}

func (d *ExcelSheet) generateSqliteDB(drivers *Drivers, storage *Storage) error {

	Parser := protoparse.Parser{}
	//加载并解析 proto文件,得到一组 FileDescriptor
	desCs, err := Parser.ParseFiles(d.rootProtoFilename)
	if err != nil {
		return errors.Errorf("GenerateSqliteDB error 生成失败解析proto：%v, %v", d.rootProtoFilename, err)
	}

	fileDescriptor := desCs[0]

	var msgDesc *desc.MessageDescriptor

	for _, v := range fileDescriptor.GetMessageTypes() {
		if d.ProtoMessageName == v.GetName() {
			msgDesc = v
		}
	}

	if msgDesc == nil {
		return errors.Errorf("Proto数据中没有找到messageName 消息名称名称：%v", d.ProtoMessageName)
	}

	curBackupDriverName := "sqlite3_backup_" + d.XlsxNameNoExt + "_" + d.Name

	driver := drivers.getOrCreate(curBackupDriverName)
	defer func() {
		driver.driverCons = nil
	}()
	driver.driverCons = nil

	srcDB, errCreateSrcDB := sql.Open(curBackupDriverName, ":memory:")

	if errCreateSrcDB != nil {
		return errors.Errorf("数据库开启失败 err %v", errCreateSrcDB)
	}
	defer srcDB.Close()

	errPing := srcDB.Ping()
	if errPing != nil {
		return errors.Errorf("Failed to connect to the source database:%v", errPing)
	}

	titleRow := d.sheet.Rows[1]
	defaultRow := d.sheet.Rows[4]
	typeRow := d.sheet.Rows[2]

	if len(typeRow.Cells) != len(titleRow.Cells) {
		return errors.Errorf("表格数据类型和标题数量不一致 表名：%s", d.FullName)
	}

	if err := createMemoryDBTable(srcDB); err != nil {
		return errors.Wrapf(err, "数据库创建表结构失败")
	}

	for j := starReadLine; j < len(d.sheet.Rows); j++ {
		msg := dynamic.NewMessage(msgDesc)

		var keyStr string

		curRow := d.sheet.Rows[j]

		isExistKey := false

		for k := 0; k < len(titleRow.Cells); k++ {

			title := titleRow.Cells[k].String()
			strType := strings.ToLower(typeRow.Cells[k].String())

			if title == "" {
				continue
			}

			if strings.ToLower(title) == "key" || strings.ToLower(title) == "id" {
				isExistKey = true
			}

			cellStr := ""

			if k < len(curRow.Cells) {
				cellStr = curRow.Cells[k].String()
			}

			for _, fieldDesc := range msgDesc.GetFields() {
				fieldName := fieldDesc.GetName()

				if strings.ToLower(fieldName) == strings.ToLower(title) {

					if strings.TrimSpace(cellStr) == "" {
						//为空之后直接去拿默认值,区别Constant 判定是存在key值得表

						if k < len(defaultRow.Cells) && strings.TrimSpace(defaultRow.Cells[k].String()) != "" {
							cellStr = defaultRow.Cells[k].String()
						} else {
							continue
						}
					}

					if strings.HasPrefix(cellStr, "**") {
						continue
					}

					if fieldDesc.GetType().String() == "TYPE_INT32" {

						if fieldDesc.IsRepeated() {

							if len(strings.Split(strType, "_")) == 2 && strings.Split(strType, "_")[1] == "list" {
								valueVec := strings.Split(cellStr, ",")

								for _, value := range valueVec {

									value = strings.TrimSpace(value)

									if value == "" {
										continue
									}

									valueInt, err := strconv.ParseInt(value, 10, 32)

									if err != nil {
										return errors.Errorf("表名：%v_%v title:%v type: %v 行数:%d,列数：%d 对应INT_LIST数据转换失败：%v 对应装换数值 %v ERR:%v", d.XlsxName, d.Name, title, strType, j+1, k+1, cellStr, value, err)
									}

									msg.AddRepeatedFieldByName(fieldDesc.GetName(), int32(valueInt))
								}
							} else {
								value, err := strconv.Atoi(cellStr)

								if err != nil {
									return errors.Errorf("表名：%v_%v 行数:%d,列数：%d 类型：%v 对应INT数据转换失败：%v ERR:%v", d.XlsxName, d.Name, j+1, k+1, strType, cellStr, err)
								}

								msg.AddRepeatedFieldByName(fieldDesc.GetName(), int32(value))
							}
						} else {

							value, err := strconv.Atoi(cellStr)

							if err != nil {
								return errors.Errorf("表名：%v_%v 行数:%d,列数：%d 对应INT数据转换失败：%v ERR:%v", d.XlsxName, d.Name, j+1, k+1, cellStr, err)
							}

							msg.SetFieldByName(fieldDesc.GetName(), int32(value))
						}

						if strings.ToLower(title) == "id" {
							keyStr = cellStr
						}
					}

					if fieldDesc.GetType().String() == "TYPE_STRING" {

						if fieldDesc.IsRepeated() {
							if len(strings.Split(strType, "_")) == 2 && strings.Split(strType, "_")[1] == "list" {
								valueVec := strings.Split(cellStr, ",")
								for _, value := range valueVec {
									msg.AddRepeatedFieldByName(fieldDesc.GetName(), value)
								}
							} else {
								msg.AddRepeatedFieldByName(fieldDesc.GetName(), cellStr)
							}
						} else {
							msg.SetFieldByName(fieldDesc.GetName(), cellStr)
						}

						if strings.ToLower(title) == "key" || strings.ToLower(title) == "id" {
							keyStr = cellStr
						}
					}

					if fieldDesc.GetType().String() == "TYPE_BOOL" {

						value, err := strconv.ParseBool(cellStr)

						if err != nil {
							return errors.Errorf("行数:%d,列数：%d 对应BOOL数据转换失败：%v ERR:%v", j, k, cellStr, err)
						}

						if fieldDesc.IsRepeated() {
							//布尔不支持重复
						} else {
							msg.SetFieldByName(fieldDesc.GetName(), value)
						}
					}

					if fieldDesc.GetType().String() == "TYPE_FLOAT" {
						if fieldDesc.IsRepeated() {

							if len(strings.Split(strType, "_")) == 2 && strings.Split(strType, "_")[1] == "list" {
								valueVec := strings.Split(cellStr, ",")
								for _, value := range valueVec {
									value = strings.TrimSpace(value)

									if value == "" {
										continue
									}

									valueFloat, err := strconv.ParseFloat(value, 32)

									if err != nil {
										return errors.Errorf("表名：%v_%v 行数:%d,列数：%d 对应FLOAT数据转换失败：%v ERR:%v", d.XlsxName, d.Name, j+1, k+1, cellStr, err)
									}

									msg.AddRepeatedFieldByName(fieldDesc.GetName(), float32(valueFloat))
								}
							} else {
								value, err := strconv.ParseFloat(cellStr, 32)

								if err != nil {
									return errors.Errorf("表名：%v_%v 行数:%d,列数：%d 对应FLOAT数据转换失败：%v ERR:%v", d.XlsxName, d.Name, j+1, k+1, cellStr, err)
								}

								msg.AddRepeatedFieldByName(fieldDesc.GetName(), float32(value))
							}
						} else {
							cellStr = strings.TrimSpace(cellStr)

							if cellStr == "" {
								continue
							}

							value, err := strconv.ParseFloat(cellStr, 32)

							if err != nil {
								return errors.Errorf("表格数据中FLOAT字段类型错误 %v", err)
							}

							msg.SetFieldByName(fieldDesc.GetName(), float32(value))
						}
					}
					break
				}
			}
		}

		//存在对应的ID 或者key title 字段
		if isExistKey {
			if keyStr == "" {
				continue
			}
		} else {
			//不存在对应的ID 或者key title 字段 基本上就是常量表 拼1个key
			if j == starReadLine {
				keyStr = "1"
			} else {
				continue
			}
		}

		errSaveDB := saveToMemoryDB(srcDB, keyStr, msg)

		if errSaveDB != nil {
			return errors.Errorf("数据库存盘失败 err %v path:%v sheetName:%v", errSaveDB, d.XlsxName, d.Name)
		}
	}

	destDb, errCreateDest := sql.Open(curBackupDriverName, d.rootSqliteFilename)

	if errCreateDest != nil {
		return errors.Errorf("数据库存盘失败，sql.Open err %v", errCreateDest)
	}
	defer destDb.Close()

	errPingDest := destDb.Ping()
	if errPingDest != nil {
		return errors.Errorf("Failed to connect to the destination database:%v", errPingDest)
	}

	if errSaveToDB := saveToDB(driver.driverCons, destDb, srcDB, d.XlsxNameNoExt, d.Name); errSaveToDB != nil {
		return errors.Errorf("数据库存盘失败 err %v", errSaveToDB)
	}

	dbBytes, err := os.ReadFile(d.rootSqliteFilename)
	if err != nil {
		return errors.Wrapf(err, "读取db文件失败")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = storage.Put(ctx, d.SqliteFilename, dbBytes); err != nil {
		return errors.Wrapf(err, "上传db文件失败")
	}

	return nil
}

func createMemoryDBTable(srcDB *sql.DB) error {

	if srcDB == nil {
		return errors.Errorf("数据库未开启")
	}

	createTableStr := "CREATE TABLE IF NOT EXISTS data  (id string PRIMARY KEY,data byte[]);"

	_, errCreate := srcDB.Exec(createTableStr)

	if errCreate != nil {
		return errors.Errorf("内存数据库创建表失败，err %v", errCreate)
	}

	return nil
}

func saveToMemoryDB(srcDB *sql.DB, keyStr string, m *dynamic.Message) error {

	var errInsert error

	//if isMarshal {
	dataMsg, errMarshal := m.Marshal()

	if errMarshal != nil {
		return errors.Errorf("Marshal err %v", errMarshal)
	}

	_, errInsert = srcDB.Exec("INSERT INTO data (id, data) VALUES (?, ?)", keyStr, dataMsg)
	//} else {
	//	_, errInsert = srcDB.Exec("INSERT INTO data (id, data) VALUES (?, ?)", keyStr, m.String())
	//}

	if errInsert != nil {
		return errors.Errorf("内存数据库插入keyStr数据失败,id:%s,err %v data:%v", keyStr, errInsert, m.String())
	}

	return nil
}

func saveToDB(driverConns []*sqlite3.SQLiteConn, destDb *sql.DB, srcDB *sql.DB, tableName string, sheetName string) error {

	if srcDB == nil {
		return errors.Errorf("保存内存数据库")
	}

	if len(driverConns) != 2 {
		return errors.Errorf("Expected 2 driver connections, but found %v.", len(driverConns))
	}

	//开始进行备份
	srcDbDriverConn := driverConns[0]

	if srcDbDriverConn == nil {
		return errors.Errorf("The source database driver connection is nil.")
	}

	destDbDriverConn := driverConns[1]

	if destDbDriverConn == nil {
		return errors.Errorf("The destination database driver connection is nil.")
	}

	// Prepare to perform the backup.
	backup, err := destDbDriverConn.Backup("main", srcDbDriverConn, "main")
	if err != nil {
		return errors.Errorf("Failed to initialize the backup:%v", err)
	}

	// Allow the initial page count and remaining values to be retrieved.
	// According to <https://www.sqlite.org/c3ref/backup_finish.html>, the page count and remaining values are "... only updated by sqlite3_backup_step()."
	isDone, err := backup.Step(0)
	if err != nil {
		return errors.Errorf("Unable to perform an initial 0-page backup step: %v", err)
	}
	if isDone {
		fmt.Println("Backup is unexpectedly done.")
	}

	// Check that the page count and remaining values are reasonable.
	initialPageCount := backup.PageCount()
	if initialPageCount <= 0 {
		fmt.Println("Unexpected initial page count value:", initialPageCount, tableName, sheetName)
		return nil
	}
	initialRemaining := backup.Remaining()
	if initialRemaining <= 0 {
		return errors.Errorf("Unexpected initial remaining value: %v", initialRemaining)
	}
	if initialRemaining != initialPageCount {
		return errors.Errorf("Initial remaining value differs from the initial page count value; remaining: %v; page count: %v", initialRemaining, initialPageCount)
	}

	// Perform the backup.
	if false {
		var startTime = time.Now().Unix()

		// Test backing-up using a page-by-page approach.
		var latestRemaining = initialRemaining
		for {
			// Perform the backup step.
			isDone, err = backup.Step(1)
			if err != nil {
				return errors.Errorf("Failed to perform a backup step:%v", err)
			}

			// The page count should remain unchanged from its initial value.
			currentPageCount := backup.PageCount()
			if currentPageCount != initialPageCount {
				return errors.Errorf("Current page count differs from the initial page count; initial page count: %v; current page count: %v", initialPageCount, currentPageCount)
			}

			// There should now be one less page remaining.
			currentRemaining := backup.Remaining()
			expectedRemaining := latestRemaining - 1
			if currentRemaining != expectedRemaining {
				return errors.Errorf("Unexpected remaining value; expected remaining value: %v; actual remaining value: %v", expectedRemaining, currentRemaining)
			}
			latestRemaining = currentRemaining

			if isDone {
				break
			}

			// Limit the runtime of the backup attempt.
			if (time.Now().Unix() - startTime) > 150 {
				return errors.Errorf("Backup is taking longer than expected.")
			}
		}
	} else {
		// Test the copying of all remaining pages.
		isDone, err = backup.Step(-1)
		if err != nil {
			return errors.Errorf("Failed to perform a backup step:%v", err)
		}
		if !isDone {
			return errors.Errorf("Backup is unexpectedly not done.")
		}
	}

	// Finish the backup.
	err = backup.Finish()
	if err != nil {
		return errors.Errorf("Failed to finish backup:%v", err)
	}

	return nil
}
