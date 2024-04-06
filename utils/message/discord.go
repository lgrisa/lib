package message

import (
	"context"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"math/rand"
	"time"
)

//StartConfig.Notify.DiscordWebhookUrl = "https://discord.com/api/webhooks/1075323981193826354/rCJrCgDxYIV3E-gpuhh6F8zh8smCnev9Tguil9flnMaI2fVMNTwbp2fYEh0yAwcWsDIX"
//StartConfig.Notify.DiscordRobotThread = "1161959914185429053"

func sendDiscord(discordWebhookUrl, discordRobotThread, title, message string) {

	client, err := webhook.NewWithURL(discordWebhookUrl)
	if err != nil {
		log.Errorf("init webhook fail failed,url:%v err: %v", discordWebhookUrl, err)
		return
	}
	defer client.Close(context.Background())

	eb := discord.NewEmbedBuilder()
	eb.SetTitlef(title)
	eb.SetDescriptionf(message)
	eb.SetTimestamp(time.Now())
	eb.SetColor(rand.Intn(0xffffff + 1))

	b := discord.NewWebhookMessageCreateBuilder()
	b.AddEmbeds(eb.Build())

	if discordRobotThread != "" {
		threadID, _ := snowflake.Parse(discordRobotThread)
		_, err = client.CreateMessageInThread(b.Build(), threadID)
	} else {
		_, err = client.CreateMessage(b.Build())
	}

	if err != nil {
		log.Errorf("webhook.CreateMessage failed, err: %v", err)
	}
}

func SendDiscordNoGoroutines(discordWebhookUrl, discordRobotThread, title, message string) {

	log.Errorf("SendDiscord message:%v", message)

	if discordWebhookUrl == "" {
		return
	}

	client, err := webhook.NewWithURL(discordWebhookUrl)
	if err != nil {
		log.Errorf("init webhook fail failed,url:%v err: %v", discordWebhookUrl, err)
		return
	}
	defer client.Close(context.Background())

	eb := discord.NewEmbedBuilder()

	eb.SetTitlef(title)
	eb.SetDescriptionf(message)
	eb.SetTimestamp(time.Now())
	eb.SetColor(rand.Intn(0xffffff + 1))

	b := discord.NewWebhookMessageCreateBuilder()
	b.AddEmbeds(eb.Build())

	if discordRobotThread != "" {
		threadID, _ := snowflake.Parse(discordRobotThread)
		_, err = client.CreateMessageInThread(b.Build(), threadID)
	} else {
		_, err = client.CreateMessage(b.Build())
	}

	if err != nil {
		log.Errorf("webhook.CreateMessage failed, err: %v", err)
	}
}
