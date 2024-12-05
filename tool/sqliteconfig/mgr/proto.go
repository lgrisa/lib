package mgr

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

func GetMessageName(filenameOnly string, sheetName string) string {
	return "Confpb" + filenameOnly + sheetName
}

func (d *ExcelSheet) CompileCs() error {

	if err := os.WriteFile(d.rootCsDirProto, d.protoBytes, os.ModePerm); err != nil {
		return errors.Wrapf(err, "写入proto文件失败")
	}

	name := "protoc"
	arg := []string{"--csharp_out=" + d.rootCsDir, d.rootCsDirProto}

	cmd := exec.Command(name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()

	if err != nil {
		return errors.Wrapf(err, "执行命令失败，Name: %s, arg: %v, stdout: %s, stderr: %s", name, arg, stdout.String(), stderr.String())
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if len(outStr) > 0 {
		fmt.Println(outStr)
	}

	if len(errStr) > 0 {
		return errors.Errorf(errStr)
	}

	return nil
}

const ProtoPrefix = "syntax = \"proto3\";\n\npackage confproto;\n\noption go_package=\"conf\";"

func (d *ExcelSheet) newProtoBytes() []byte {
	//消息头部
	b := bytes.Buffer{}
	b.WriteString(ProtoPrefix)
	b.WriteString("\nmessage " + d.ProtoMessageName + " {\n")

	for _, f := range d.FieldArray {

		//如果首字母是小写，将首字母转为大写
		if f.Name[0] >= 'a' && f.Name[0] <= 'z' {
			f.Name = strings.ToUpper(f.Name[:1]) + f.Name[1:]
		}

		if f.IsRepeated {
			b.WriteString(fmt.Sprintf("    repeated %v %v = %v;", f.ProtoType, f.Name, f.ProtoId))
		} else {
			b.WriteString(fmt.Sprintf("    %v %v = %v;", f.ProtoType, f.Name, f.ProtoId))
		}
		b.WriteString("\n")
	}

	b.WriteString("}\n")

	return b.Bytes()
}

func getProtoType(titleType string) (string, bool) {
	protoType := "string"

	if titleType == "bool" {
		protoType = "bool"
	}

	if titleType == "float" {
		protoType = "float"
	}

	if titleType == "int" {
		protoType = "int32"
	}

	typeArray := strings.Split(titleType, "_")

	if len(typeArray) == 2 && typeArray[1] == "list" {

		repeatedType := "string"

		if typeArray[0] == "bool" {
			repeatedType = "bool"
		}

		if typeArray[0] == "float" {
			repeatedType = "float"
		}

		if typeArray[0] == "int" {
			repeatedType = "int32"
		}

		return repeatedType, true
	}

	return protoType, false
}
