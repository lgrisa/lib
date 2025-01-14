package utils

import (
	"fmt"
	"github.com/lgrisa/lib/utils/logutil"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func GetPreviousBetweenCommit(tag string) string {
	tagList := strings.Split(tag, ".")
	if len(tagList) < 2 {
		logutil.LogTraceF("tag format err: %s not get commit\n", tag)
		return ""
	}

	lastTagNum := tagList[len(tagList)-1]

	lastTagNumInt, err := strconv.Atoi(lastTagNum)
	if err != nil {
		logutil.LogTraceF("tag转换数字失败: %s, err: %s\n", lastTagNum, err)
		return ""
	}

	if lastTagNumInt < 1 {
		logutil.LogTraceF("tag num < 1: %s\n", lastTagNum)
		return ""
	}

	lastTagName := strings.Join(tagList[0:len(tagList)-1], ".") + "." + strconv.Itoa(lastTagNumInt-1)

	//检查是否存在tag
	output, err := RunCommandGetOutPut(fmt.Sprintf("git tag -l %s", lastTagName))

	if err != nil {
		logutil.LogTraceF("获取tag失败: %s, err: %s\n", lastTagName, err)
		return ""
	}

	if string(output) == "" {
		logutil.LogTraceF("tag不存在: %s\n", lastTagName)
		return ""
	}

	//%H   提交对象（commit）的完整哈希字串
	//%h    提交对象的简短哈希字串
	//%T    树对象（tree）的完整哈希字串
	//%t    树对象的简短哈希字串
	//%P    父对象（parent）的完整哈希字串
	//%p    父对象的简短哈希字串
	//%an   作者（author）的名字
	//%ae   作者的电子邮件地址
	//%ad   作者修订日期（可以用 -date= 选项定制格式）
	//%ar   作者修订日期，按多久以前的方式显示
	//%cn   提交者(committer)的名字
	//%ce   提交者的电子邮件地址
	//%cd   提交日期
	//%cr   提交日期，按多久以前的方式显示
	//%s    提交说明

	//获取tag之间的commit
	output, err = RunCommandGetOutPut(fmt.Sprintf("git logutil %s...%s --pretty=format:\"%%s\"", lastTagName, tag))

	if err != nil {
		logutil.LogErrorF("获取tag之间的commit失败: %s, err: %s\n", lastTagName, err)
		return ""
	}

	//美化下返回值，每行中间加一行空白行
	returnStr := "*" + string(output)

	return strings.ReplaceAll(returnStr, "\n", "\n\n*")
}

func GetTagCommit(gitPath string, tag string) (string, error) {
	//获取tag之间的commit
	gitCmd := fmt.Sprintf("git logutil %s --pretty=format:\"%%s\"", tag)

	if gitPath != "" {
		gitCmd = fmt.Sprintf("cd %s && %s", gitPath, gitCmd)
	}

	output, err := RunCommandGetOutPut(gitCmd)

	if err != nil {
		return "", errors.Errorf("get tag commit failed: %s, err: %s", tag, err)
	}

	//美化下返回值，每行中间加一行空白行
	returnStr := "*" + string(output)

	if gitPath != "" {
		_, _ = RunCommandGetOutPut("cd -")
	}

	return strings.ReplaceAll(returnStr, "\n", "\n\n*"), nil
}
