package mgr

import (
	"bytes"
	"fmt"
	"github.com/lgrisa/lib/utils"
	"github.com/lgrisa/lib/utils/compress"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"text/template"
	"time"
	"unicode"
)

func (m *Manager) handleGenTypeScript(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("invalid file size", fileSize, "  val :", val)
		writeErrMsg(w, "invalid size")
		return
	}
	fileSizeInt64 = val

	fmt.Println("handleGenTypeScript", fileMd5, fileSize)

	//重新生成
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("io.ReadAll(r.Body) fail: ", err)
		writeErrMsg(w, "io.ReadAll(r.Body) fail: "+err.Error())
		return
	}

	dataMd5 := utils.Md5String(data)
	dataSize := int64(len(data))

	if fileMd5 != dataMd5 || fileSizeInt64 != dataSize {
		fmt.Println(fmt.Sprintf("fileMd5(%v) != DataMd5(%v) || fileSizeInt64(%v) != dataSize(%v)", fileMd5, dataMd5, fileSizeInt64, dataSize))
		writeErrMsg(w, fmt.Sprintf("header和body的数据不一致, headerMd5: %v, bodyMd5: %v, headerSize: %v, bodySize: %v", fileMd5, dataMd5, fileSizeInt64, dataSize))
		return
	}

	if err := m.funcIdMap(func(idMap *MessageIdGen, drivers *Drivers) bool {

		excelZip := newExcelZip(data, dataMd5, m.root, packVersion, idMap, drivers, m.storage)
		//defer os.RemoveAll(excelZip.rootCsDir)

		if err := excelZip.doGenerate(false, false, false); err != nil {
			errMsg := fmt.Sprintf("生成Cs文件失败, %v", err)
			fmt.Println(errMsg)
			writeErrMsg(w, errMsg)
			return false
		}

		protoBytes, err := excelZip.genTypeScriptProto()
		if err != nil {
			errMsg := fmt.Sprintf("genTypeScriptProto fail, %v", err)
			fmt.Println(errMsg)
			writeErrMsg(w, errMsg)
			return false
		}

		pbjs, pbts, err := compileTypeScriptProtobuf(protoBytes)
		if err != nil {
			errMsg := fmt.Sprintf("compileTypeScriptProtobuf fail, %v", err)
			fmt.Println(errMsg)
			writeErrMsg(w, errMsg)
			return false
		}

		configBytes, err := excelZip.genTypeScriptConfigs()
		if err != nil {
			errMsg := fmt.Sprintf("genTypeScriptConfigs fail, %v", err)
			fmt.Println(errMsg)
			writeErrMsg(w, errMsg)
			return false
		}

		fileBytes := make(map[string][]byte)
		fileBytes["Protobuf/confpb.proto"] = protoBytes
		fileBytes["Protobuf/confpb.js"] = pbjs
		fileBytes["Protobuf/confpb.d.ts"] = pbts
		fileBytes["Logic/Config/ConfigsGen.ts"] = configBytes

		body, err := compress.PackZipData(fileBytes)
		if err != nil {
			errMsg := fmt.Sprintf("PackData fail, %v", err)
			fmt.Println(errMsg)
			writeErrMsg(w, errMsg)
			return false
		}

		confJson := &ConfigCsJson{}
		confJson.FileMd5 = dataMd5
		confJson.FileSize = dataSize
		confJson.CsBody = body

		writeCsJson(w, confJson, "使用生成的版本:"+fileMd5)
		return true
	}); err != nil {
		errMsg := fmt.Sprintf("funcIdMap fail, %v", err)
		fmt.Println(errMsg)
		writeErrMsg(w, errMsg)
		return
	}
}

func compileTypeScriptProtobuf(protoBytes []byte) ([]byte, []byte, error) {

	basePath := fmt.Sprintf("/tmp/%v", time.Now().UnixNano())
	os.MkdirAll(basePath, os.ModePerm)
	defer os.RemoveAll(basePath)

	// 写入proto文件
	protoPath := fmt.Sprintf("%v/confpb.proto", basePath)
	if err := os.WriteFile(protoPath, protoBytes, os.ModePerm); err != nil {
		return nil, nil, errors.Wrapf(err, "写入proto文件失败："+protoPath)
	}

	// 编译proto文件
	cmd := exec.Command("pbjs", "--dependency", "protobuf", "-t", "static-module", "-w", "commonjs", "-o", "confpb.js", "confpb.proto")
	cmd.Dir = basePath

	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Println(string(out))
	}
	if err != nil {
		return nil, nil, errors.Wrapf(err, "compile pb proto pbjs fail, workdir: %v", protoPath)
	}

	cmd = exec.Command("pbts", "-o", "confpb.d.ts", "confpb.js")
	cmd.Dir = basePath

	out, err = cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Println(string(out))
	}
	if err != nil {
		return nil, nil, errors.Wrapf(err, "compile pb proto pbts fail, workdir: %v", protoPath)
	}

	protoJsPath := fmt.Sprintf("%v/confpb.js", basePath)
	protoDtsPath := fmt.Sprintf("%v/confpb.d.ts", basePath)

	protoJsBytes, err := os.ReadFile(protoJsPath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "读取protoJs文件失败："+protoJsPath)
	}

	protoDtsBytes, err := os.ReadFile(protoDtsPath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "读取protoDts文件失败："+protoDtsPath)
	}

	return protoJsBytes, protoDtsBytes, nil
}

func (d *excel_zip) genTypeScriptProto() ([]byte, error) {
	var protoBuf bytes.Buffer

	//proto文件头部
	protoBuf.WriteString(ProtoPrefix)

	var ks []string
	for k := range d.fileMap {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	for _, k := range ks {
		f := d.fileMap[k]

		var kss []string
		for kk := range f.SheetMap {
			kss = append(kss, kk)
		}
		sort.Strings(kss)

		for _, kk := range kss {
			s := f.SheetMap[kk]
			protoFilenameData, err := os.ReadFile(s.rootProtoFilename)
			if err != nil {
				return nil, errors.Wrapf(err, "读取proto文件失败, %v", s.rootProtoFilename)
			}

			//处理掉开头
			protoFilenameData = bytes.TrimPrefix(protoFilenameData, []byte(ProtoPrefix))

			//合并proto文件
			protoBuf.Write(protoFilenameData)

			//写入换行
			protoBuf.WriteString("\n")
		}
	}

	return protoBuf.Bytes(), nil
}

const typeScriptTemplate = `
import { $ref } from 'puerts';
import * as UE from 'ue'
import { confproto as pb } from "../../Protobuf/confpb"
import DB from "../Utils/sql/sql"
import bytes from '../Utils/bytes';

class ConfigDatas {

	TotalDataCount: number = 0

    public IsReady(): boolean {
        return this.dataPathMap.size >= this.TotalDataCount
    }

    private dataPathMap: Map<string, string> = new Map()
    public SetDataPath(key: string, path: string) {
        this.dataPathMap.set(key, path)
    }

    private loadData(key: string): Uint8Array {
        let path = this.dataPathMap.get(key)
        if (!path) {
            console.error("ConfigDatas.loadData path not found", key)
            return null
        }

        let fileData = UE.NewArray(UE.BuiltinByte);
        if (!UE.BaseFilesDownloader.LoadFileToArray(path, $ref(fileData))) {
            console.error("文件加载失败", path)
            return;
        }

        let data = bytes.FromByteArray(fileData)
        console.log("文件本地加载成功", path, data.length)

        return data
    }

    private loadDb(key: string): DB {
        let data = this.loadData(key)
        if (!data) {
            return null
        }

        return DB.Create(data)
    }

{{range .Configs}}
    private {{firstLower .Name}}Config: {{.Name}}Config

    private get{{.Name}}Config(): {{.Name}}Config {
        if (!this.{{firstLower .Name}}Config) {
            this.{{firstLower .Name}}Config = {{.Name}}Config.New(this.loadDb("{{.ProtoName}}"))
        }
        return this.{{firstLower .Name}}Config
    }

    public Get{{.Name}}Array(): {{.Name}}[] {
        let config = this.get{{.Name}}Config()
        return config.Array
    }

    public Get{{.Name}}(id: string): {{.Name}} {
        let config = this.get{{.Name}}Config()
        return config.Map.get(id)
    }
{{end}}
}

{{range .Configs}}
class {{.Name}}Config {
    public static New(db: DB): {{.Name}}Config {
        let config = new {{.Name}}Config()
        let ids = db.LoadIds()
        for (let id of ids) {
            let data = new {{.Name}}()
            data.db = db
            data.id = id
            config.Array.push(data)
            config.Map.set(id, data)
        }
        return config;
    }

    public Array: {{.Name}}[] = []
    public Map: Map<string, {{.Name}}> = new Map()
}

class {{.Name}} {
    db: DB

    id: string

    data: pb.{{.ProtoName}}

    public GetData(): pb.{{.ProtoName}} {
        if (!this.data) {
            this.data = pb.{{.ProtoName}}.decode(this.db.LoadData(this.id))
        }
        return this.data
    }
}
{{end}}

export default ConfigDatas;
`

func (d *excel_zip) genTypeScriptConfigs() ([]byte, error) {
	buf := &bytes.Buffer{}

	type config_data struct {

		// HeroData
		Name string

		// HeroDataProto
		ProtoName string
	}

	type typeScriptStruct struct {
		Configs []*config_data
	}

	v := &typeScriptStruct{}

	var ks []string
	for k := range d.fileMap {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	for _, k := range ks {
		f := d.fileMap[k]

		var kss []string
		for kk := range f.SheetMap {
			kss = append(kss, kk)
		}
		sort.Strings(kss)

		for _, kk := range kss {
			s := f.SheetMap[kk]

			v.Configs = append(v.Configs, &config_data{
				Name:      HumpName(s.XlsxNameNoExt + "_" + s.Name),
				ProtoName: s.ProtoMessageName,
			})
		}
	}

	temp := template.Must(template.New("configs").
		Funcs(template.FuncMap{
			"hump":           HumpName,
			"firstLowerHump": firstLowerHumpName,
			"firstLower":     firstLowerCase,
		}).Parse(typeScriptTemplate))

	if err := temp.Execute(buf, v); err != nil {
		return nil, errors.Wrapf(err, "模板生成失败")
	}

	return buf.Bytes(), nil
}

func firstLowerCase(s string) string {
	var buf bytes.Buffer
	for i, r := range s {
		if i == 0 {
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func firstLowerHumpName(s string) string {
	return firstLowerCase(HumpName(s))
}
