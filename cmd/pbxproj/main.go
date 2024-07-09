package main

import (
	"fmt"
	"github.com/lgrisa/lib/pbxproj"
	"github.com/lgrisa/lib/utils"
	"os"
	"strings"
)

func main() {
	utils.InitLog()

	projPath, isFound := utils.FindProjectPath("")

	utils.LogDebugF("projPath:%v", projPath)

	if !isFound {
		utils.LogDebugF("not found project.pbxproj")
		return
	}

	// parse pbxproj
	proj, err := pbxproj.NewPbxproj(projPath)

	if err != nil {
		utils.LogErrorF("parse pbxproj err:%v", err)
		return
	}

	rootObjectUUid, err := proj.GetJson().Get("rootObject").String()

	if err != nil {
		utils.LogDebugF("get rootObject err:%v", err)
		return
	}

	utils.LogDebugF("rootObject:%v", rootObjectUUid)

	pBXProject, isFound := proj.ProjectSection[rootObjectUUid]
	if !isFound {
		utils.LogDebugF("not found PBXProject:%v", rootObjectUUid)
		return
	}

	var shouldDeleteIds []string

	if unityFrameworkTarget, isFound := proj.NativeTargets[pBXProject.TargetsUnityFramework]; isFound {

		utils.LogDebugF("unityFrameworkTarget.BuildPhasesFrameworks:%v", unityFrameworkTarget.BuildPhasesFrameworks)

		if buildPhase, isFound := proj.FrameworkBuildPhase[unityFrameworkTarget.BuildPhasesFrameworks]; isFound {
			for _, file := range buildPhase.Files {
				shouldDeleteIds = append(shouldDeleteIds, file)
			}
		}
	}

	//删除Bugly.framework
	if err = rmBuglyFramework(projPath, buglyStr, proj.Lines, shouldDeleteIds); err != nil {
		utils.LogDebugF("rmBuglyFramework err:%v", err)
	}
}

const buglyStr = "/* Bugly.framework in Frameworks */"

func rmBuglyFramework(filePath string, deleteTagStr string, lines []string, filesIds []string) error {

	var result []string

	for _, line := range lines {

		isContain := false
		for _, filesId := range filesIds {
			if strings.Contains(line, deleteTagStr) && strings.Contains(line, filesId) {
				//找到了Bugly.framework
				//删除这一行
				isContain = true
				break
			}
		}

		if isContain {
			fmt.Print("delete line:", line)
			continue
		}

		//将读取到的每一行数据存放到切片中
		result = append(result, line)
	}

	err := os.RemoveAll(filePath)

	if err != nil {
		return err
	}

	//创建一个新的文件
	newFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	//将切片中的数据写入到新文件中
	for _, v := range result {
		_, err = newFile.WriteString(v)
		if err != nil {
			return err
		}
	}

	//关闭文件
	return newFile.Close()
}
