package main

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/lgrisa/lib/utils"
	"github.com/lgrisa/lib/utils/log"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

//https://discohook.org/

const (
	DiscordUrl  = "https://discord.com/api/webhooks/1075323981193826354/rCJrCgDxYIV3E-gpuhh6F8zh8smCnev9Tguil9flnMaI2fVMNTwbp2fYEh0yAwcWsDIX"
	RobotThread = "1075322890276319232"

	BuildDir   = "gitlab-runner-build"
	QRCodePath = "qrCode"

	uploadName           = "uploader"
	uploadPasswd         = "RM1k-ys["
	uploadHost           = "192.168.0.200"
	uploadTargetRootPath = "web"
	uploadTargetMidPath  = "test"

	downloadUrl = "https://nasload.soulframegame.com"
)

//goMac.exe --name=./build/send --path=./cmd/packer/main.go

type initCiStruct struct {
	tag            string
	shellName      string
	commitText     string
	activeDirPath  string
	activeDir      []os.DirEntry
	sendMessageMap map[string]*sendDisCordStruct
}

func main() {

	for i, v := range os.Args {
		log.LogTracef("args[%d]: %s\n", i, v)
	}

	output, _ := utils.RunCommandGetOutPut("pwd")
	log.LogTracef("cur dir: %s", string(output))

	initCi, err := initStruct()

	if err != nil {
		execExit(fmt.Sprintf("init struct failed: %s", err))
	}

	if err = initCi.work(); err != nil {
		execExit(fmt.Sprintf("work failed: %s", err))
	}

	execExit("")
}

func initStruct() (*initCiStruct, error) {
	for i, v := range os.Args {
		log.LogTracef("args[%d]: %s\n", i, v)
	}

	execTag := os.Args[1]

	execShell := os.Args[2]

	if strings.HasSuffix(execShell, ".sh") {
		err := utils.RunSyncCommand(fmt.Sprintf("chmod +x %s && ./%s", execShell, execShell))
		if err != nil {
			return nil, errors.Wrapf(err, "run shell failed")
		}
	} else {
		log.LogTracef("Not shell script %s", execShell)
	}

	err := os.MkdirAll(BuildDir, os.ModePerm)
	if err != nil {
		return nil, errors.New("create build dir failed: " + err.Error())
	}

	err = os.MkdirAll(QRCodePath, os.ModePerm)
	if err != nil {
		return nil, errors.New("create qrcode dir failed: " + err.Error())
	}

	activeDirPath := ""
	if len(os.Args) < 4 {
		activeDirPath = "/Users/packer/SGPackage"
	} else {
		activeDirPath = os.Args[3] //生成的文件夹
	}

	platform := getPlatformByExecTag(execTag)

	log.LogTracef("platform: %s", platform)

	if platform == "" {
		execExit(fmt.Sprintf("unkonw platform: %s", execTag))
	}

	activeDirPath = fmt.Sprintf("%s/%s", activeDirPath, platform)

	log.LogTracef("activeDirPath: %s", activeDirPath)

	activeDir, err := os.ReadDir(activeDirPath)
	if err != nil {
		return nil, errors.Wrapf(err, "read active dir failed: %s", activeDirPath)
	}

	return &initCiStruct{
		tag:            execTag,
		shellName:      execShell,
		commitText:     utils.GetPreviousBetweenCommit(execTag),
		activeDirPath:  activeDirPath,
		activeDir:      activeDir,
		sendMessageMap: make(map[string]*sendDisCordStruct),
	}, nil
}

func (i *initCiStruct) work() error {
	if err := i.packAndSendFile(); err != nil {
		return err
	}

	if err := i.sendDiscordMessage(); err != nil {
		return err
	}

	return nil
}

func (i *initCiStruct) packAndSendFile() error {

	for _, file := range i.activeDir {

		//如果是隐藏文件，不处理
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		fileName := file.Name()

		//先压缩文件夹
		if file.IsDir() {
			fileName = fmt.Sprintf("%s.zip", file.GetFilePrefix(file.Name()))

			if _, err := utils.RunCommandGetOutPut(fmt.Sprintf("cd %s && zip -r %s %s && cd -", i.activeDirPath, fileName, file.Name())); err != nil {
				execExit(fmt.Sprintf("zip err: %s, err: %s", file.Name(), err))
			}

			log.LogTracef("zip file: %s", fileName)
		}

		//反复打包不能覆盖，对于打包出来的文件夹名称进行修饰
		newAppName := newPackerName(fileName)

		if _, err := utils.RunCommandGetOutPut(fmt.Sprintf("mv %s/%s %s/%s", i.activeDirPath, fileName, BuildDir, newAppName)); err != nil {
			execExit(fmt.Sprintf("mv err: %s, err: %s", file.Name(), err))
		}

		output, err := utils.RunCommandGetOutPut(fmt.Sprintf("ncftpput -u %s -p %s %s /%s/%s %s/%s",
			uploadName, uploadPasswd, uploadHost, uploadTargetRootPath, uploadTargetMidPath, BuildDir, newAppName))

		if err != nil {
			execExit(fmt.Sprintf("upload: %s, err: %s, output: %s", file.Name(), err, output))
		}

		pngPath := fmt.Sprintf("%s/%s", QRCodePath, fmt.Sprintf("%s.png", file.GetFilePrefix(newAppName)))

		err = qrcode.WriteFile(pngPath, qrcode.Medium, 256, pngPath)

		if err != nil {
			execExit(fmt.Sprintf("生成二维码失败: %s, err: %s", file.Name(), err))
		}

		log.LogTracef("create png: %s", pngPath)

		output, err = utils.RunCommandGetOutPut(fmt.Sprintf("ncftpput -u %s -p %s %s /%s/%s/%s %s",
			uploadName, uploadPasswd, uploadHost, uploadTargetRootPath, uploadTargetMidPath, QRCodePath, pngPath))

		if err != nil {
			execExit(fmt.Sprintf("上传二维码失败: %s, err: %s", file.Name(), err))
		}

		i.sendMessageMap[file.Name()] = &sendDisCordStruct{
			Title:   newAppName,
			Url:     fmt.Sprintf("%s/%s/%s", downloadUrl, uploadTargetMidPath, newAppName),
			PngUrl:  fmt.Sprintf("%s/%s/%s", downloadUrl, uploadTargetMidPath, pngPath),
			PngPath: pngPath,
		}
	}

	return nil
}

func (i *initCiStruct) sendDiscordMessage() error {
	client, err := webhook.NewWithURL(DiscordUrl)
	if err != nil {
		execExit(fmt.Sprintf("create discord client failed, err: %s", err))
	}
	defer client.Close(context.Background())

	for _, v := range i.sendMessageMap {
		eb := discord.NewEmbedBuilder()
		eb.SetAuthorName("gitlab-runner自动打包完成")
		eb.SetTitle(v.Title)
		eb.SetURL(v.Url)
		eb.SetColor(rand.Intn(0xffffff + 1))
		eb.SetTimestamp(time.Now())

		if i.tag != "" {
			eb.AddFields(discord.EmbedField{Name: "执行Tag", Value: i.tag})
		}

		if i.commitText != "" {
			eb.AddFields(discord.EmbedField{Name: "更新内容", Value: i.commitText})
		}

		eb.AddFields(discord.EmbedField{Name: "二维码地址", Value: v.PngUrl})

		eb.SetColor(rand.Intn(0xffffff + 1))

		b := discord.NewWebhookMessageCreateBuilder()
		b.SetUsername("Gitlab Runner")
		b.AddEmbeds(eb.Build())

		var f *os.File
		f, err = os.Open(v.PngPath)

		if err != nil {
			execExit(fmt.Sprintf("open png file failed, err: %s", err))
		}

		b.AddFile(v.Title+".png", v.PngPath, f)

		threadID, _ := snowflake.Parse(RobotThread)
		_, err = client.CreateMessageInThread(b.Build(), threadID)

		if err != nil {
			execExit(fmt.Sprintf("send message failed, err: %s", err))
		}
	}

	return nil
}

type sendDisCordStruct struct {
	Title   string
	Url     string
	PngUrl  string
	PngPath string
}

func execExit(errMsg string) {
	if errMsg != "" {
		log.LogErrorf(errMsg)
		senDiscordErrMsg(errMsg)
	}

	_ = os.RemoveAll(BuildDir)
	_ = os.RemoveAll(QRCodePath)

	os.Exit(0)
}

func senDiscordErrMsg(errMsg string) {
	fmt.Println(errMsg)

	client, err := webhook.NewWithURL(DiscordUrl)
	if err != nil {
		execExit(fmt.Sprintf("create discord client failed, err: %s", err))
	}
	defer client.Close(context.Background())

	eb := discord.NewEmbedBuilder()
	eb.SetAuthorName("gitlab-runner自动打包失败")
	eb.SetTitle(errMsg)
	eb.SetColor(rand.Intn(0xffffff + 1))
	eb.SetTimestamp(time.Now())

	b := discord.NewWebhookMessageCreateBuilder()
	b.SetUsername("Gitlab Runner")
	b.AddEmbeds(eb.Build())

	threadID, _ := snowflake.Parse(RobotThread)
	_, err = client.CreateMessageInThread(b.Build(), threadID)

	if err != nil {
		execExit(fmt.Sprintf("send message failed, err: %s", err))
	}
}

func newPackerName(fileName string) string {
	fileNameAll := path.Base(fileName)
	fileSuffix := path.Ext(fileNameAll)
	filePrefix := fileNameAll[0 : len(fileNameAll)-len(fileSuffix)]

	return fmt.Sprintf("%s_%s%s", filePrefix, time.Now().Format("20060102150405"), fileSuffix)
}

func getPlatformByExecTag(execTag string) string {
	lowerStr := strings.ToLower(execTag)

	if strings.Contains(lowerStr, "android") {
		return "Android"
	}

	if strings.Contains(lowerStr, "ios") {
		return "IOS"
	}

	return ""
}
