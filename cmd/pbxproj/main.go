package main

import (
	"fmt"
	"github.com/lgrisa/library/pbxproj"
	"github.com/lgrisa/library/utils"
	"os"
	"strings"
)

// only for Mac

func main() {
	projPath, isFound := utils.FindProjectPath("")

	fmt.Println(projPath, isFound)

	if projPath, isFound = utils.FindProjectPath(""); !isFound {
		fmt.Println("not found project.pbxproj")
		return
	}

	// parse pbxproj
	proj, err := pbxproj.NewPbxproj(projPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	str, err := proj.GetJson().Get("rootObject").String()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("rootObject:", str)

	pBXProject, isFound := proj.ProjectSection[str]
	if !isFound {
		fmt.Println("not found PBXProject:", str)
		return
	}

	fmt.Printf("unityIphoneTestsId:%v unityFrameworkId:%v \n", pBXProject.TargetsUnityIphoneTests, pBXProject.TargetsUnityFramework)

	var shouldDeleteIds []string

	if unityIphoneTarget, isFound := proj.NativeTargets[pBXProject.TargetsUnityIphoneTests]; isFound {

		fmt.Println("unityFrameworkTarget.BuildPhasesFrameworks:", unityIphoneTarget.BuildPhasesFrameworks)

		if buildPhase, isFound := proj.FrameworkBuildPhase[unityIphoneTarget.BuildPhasesFrameworks]; isFound {
			for _, file := range buildPhase.Files {
				shouldDeleteIds = append(shouldDeleteIds, file)
			}
		}
	}

	if unityFrameworkTarget, isFound := proj.NativeTargets[pBXProject.TargetsUnityFramework]; isFound {

		fmt.Println("unityFrameworkTarget.BuildPhasesFrameworks:", unityFrameworkTarget.BuildPhasesFrameworks)

		if buildPhase, isFound := proj.FrameworkBuildPhase[unityFrameworkTarget.BuildPhasesFrameworks]; isFound {
			for _, file := range buildPhase.Files {
				shouldDeleteIds = append(shouldDeleteIds, file)
			}
		}
	}

	//删除Bugly.framework
	if err = rmBuglyFramework(projPath, buglyStr, proj.Lines, shouldDeleteIds); err != nil {
		fmt.Println("rmBuglyFramework err:", err)
	}
}

const buglyStr = "/* Bugly.framework in Frameworks */"

func rmBuglyFramework(filePath string, deleteTagStr string, lines []string, filesIds []string) error {

	var result []string

	for _, line := range lines {

		isContain := false
		for _, filesId := range filesIds {
			if strings.Contains(line, buglyStr) && strings.Contains(line, filesId) {
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
