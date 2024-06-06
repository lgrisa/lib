package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

//goMac.exe --name=build/discordSender --path=./cmd/discord/main.go

var discordUrl = flag.String("discordUrl", "", "discordUrl")

var robotThread = flag.String("robotThread", "", "robotThread")

var userName = flag.String("userName", "", "userName")

var execTime = flag.String("execTime", "", "execTime")

var gitTag = flag.String("gitTag", "", "gitTag")

var gitPath = flag.String("gitPath", "", "gitPath")

var discordContent = flag.String("discordContent", "", "discordContent")

var pngPath = flag.String("pngPath", "", "pngPath")

func main() {
	flag.Parse()

	if *discordUrl == "" {
		fmt.Println("please input discordUrl")
		return
	}

	if *robotThread == "" {
		fmt.Println("please input robotThread")
		return
	}

	if *userName == "" {
		fmt.Println("please input userName")
		return
	}

	if *pngPath != "" {
		_, err := os.Stat(*pngPath)

		if err != nil {
			fmt.Printf("file not exist: %v\n", err)
			return
		}
	}

	if err := sendDiscordMessage(); err != nil {
		fmt.Printf("send discord message failed: %v\n", err)
		return
	}
}

func sendDiscordMessage() error {
	client, err := webhook.NewWithURL(*discordUrl)
	if err != nil {
		return errors.Wrapf(err, "create discord client failed")
	}
	defer client.Close(context.Background())

	b := discord.NewWebhookMessageCreateBuilder()
	b.SetUsername(fmt.Sprintf(*userName+"_%s", time.Now().Format("2006-01-02 15:04:05")))

	if *execTime != "" {
		*discordContent += fmt.Sprintf("\n**execTime**\n%s", *execTime)
	}

	if *gitTag != "" && *gitPath != "" {
		*discordContent += fmt.Sprintf("\n**执行tag**\n%s", *gitTag)

		if tagCommit, err := utils.GetTagCommit(*gitPath, *gitTag); err != nil {
			fmt.Println(err)
		} else {
			*discordContent += fmt.Sprintf("\n**tagCommit**\n%s", tagCommit)
		}
	}

	//存在换行符没有生效
	if strings.Contains(*discordContent, "\\n") {
		*discordContent = strings.ReplaceAll(*discordContent, "\\n", "\n")
	}

	if *discordContent != "" {
		b.SetContent(*discordContent)
	}

	if *pngPath != "" {
		var f *os.File
		f, err = os.Open(*pngPath)

		if err != nil {
			return errors.Wrapf(err, "open png file failed")
		}

		b.AddFile("file.png", "", f)
	}

	threadID, _ := snowflake.Parse(*robotThread)
	_, err = client.CreateMessageInThread(b.Build(), threadID)

	if err != nil {
		return errors.Wrapf(err, "send discord message failed")
	}

	return nil
}
