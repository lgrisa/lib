package mgr

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
	"os"
	"sort"
	"time"
)

type MsgToDB struct {
	MsgName   string
	FileName  string
	TableName string
	SheetName string
}

type VersionText struct {
	CellList []*MsgToDB
}

// {root}/{version}/version/version_{zip_md5}_{zip_len}.json
func getVersionName(root, version, fileMd5 string, fileSizeInt64 int) (string, string) {
	filename := fmt.Sprintf("%v/version/version_%s_%d.json", version, fileMd5, fileSizeInt64)
	return root + "/" + filename, filename
}

func getConfigJsonName(root, version, fileMd5 string, fileSizeInt64 int) (string, string) {
	filename := fmt.Sprintf("%v/version/config_%s_%d.json", version, fileMd5, fileSizeInt64)
	return root + "/" + filename, filename
}

func (d *excel_zip) writeVersionFile() error {
	rootFileName, versionName := getVersionName(d.root, d.version, d.dataMd5, len(d.data))

	v := &VersionText{}

	for _, f := range d.fileMap {
		for _, s := range f.SheetMap {
			newDBInfo := &MsgToDB{
				MsgName:   s.ProtoMessageName,
				FileName:  s.SqliteFilename,
				TableName: s.XlsxNameNoExt,
				SheetName: s.Name,
			}

			v.CellList = append(v.CellList, newDBInfo)
		}
	}

	sort.Slice(v.CellList, func(i, j int) bool {
		return v.CellList[i].MsgName < v.CellList[j].MsgName
	})

	data, err := jsoniter.Marshal(v)
	if err != nil {
		return errors.Wrapf(err, "version json.Marshal fail")
	}

	if err := os.WriteFile(rootFileName, data, os.ModePerm); err != nil {
		return errors.Wrapf(err, "excel_zip.writeVersionFile writeFile fail")
	}

	// 写入ConfigJson文件
	confJson := &ConfigJson{
		FileMd5:     d.dataMd5,
		FileSize:    int64(len(d.data)),
		VersionPath: versionName,
		VersionMd5:  utils.Md5String(data),
	}

	configData, err := jsoniter.Marshal(confJson)
	if err != nil {
		return errors.Wrapf(err, "config json.Marshal fail")
	}

	rootConfigFileName, _ := getConfigJsonName(d.root, d.version, d.dataMd5, len(d.data))

	if err := os.WriteFile(rootConfigFileName, configData, os.ModePerm); err != nil {
		return errors.Wrapf(err, "excel_zip.writeConfigJsonFile writeFile fail")
	}

	// 上传cdn
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.storage.Put(ctx, versionName, data); err != nil {
		return errors.Wrapf(err, "上传version文件失败")
	}

	return nil
}

func loadConfigJson(root, version, fileMd5 string, fileSizeInt64 int) (*ConfigJson, error) {
	rootConfigFilename, _ := getConfigJsonName(root, version, fileMd5, fileSizeInt64)
	configData, err := os.ReadFile(rootConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, errors.Wrapf(err, "loadConfigJson readFile fail, %v", rootConfigFilename)
	}

	c := &ConfigJson{}
	if err = jsoniter.Unmarshal(configData, c); err != nil {
		return nil, errors.Wrapf(err, "loadConfigJson jsoniter.Unmarshal fail, %v", rootConfigFilename)
	}

	return c, nil
}
