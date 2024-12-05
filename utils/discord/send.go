package discord

import (
	"context"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"strings"
	"time"
)

func SendDiscordNoGoroutines(message *Message) error {

	if message.DiscordWebhookUrl == "" {
		return errors.Errorf("discordWebhookUrl is empty")
	}

	client, err := webhook.NewWithURL(message.DiscordWebhookUrl)
	if err != nil {
		return errors.Errorf("init webhook fail failed,url:%v err: %v", message.DiscordWebhookUrl, err)
	}
	defer client.Close(context.Background())

	eb := discord.NewEmbedBuilder()

	if message.AuthorName != "" {
		eb.SetAuthorName(message.AuthorName)
	}

	eb.SetTitlef(message.Title)

	if message.Url != "" {
		eb.SetURL(message.Url)
	}

	eb.SetDescriptionf(message.Message)
	eb.SetTimestamp(time.Now())
	eb.SetColor(rand.Intn(0xffffff + 1))

	if message.FieldMap != nil {
		for k, v := range message.FieldMap {
			eb.AddFields(discord.EmbedField{Name: k, Value: v})
		}
	}

	b := discord.NewWebhookMessageCreateBuilder()

	//存在换行符没有生效
	if strings.Contains(message.Content, "\\n") {
		message.Content = strings.ReplaceAll(message.Content, "\\n", "\n")
	}

	b.SetContent(message.Content)

	if message.UserName != "" {
		b.SetUsername(message.UserName)
	}

	b.AddEmbeds(eb.Build())

	if message.DiscordFile != nil && message.DiscordFile.FilePath != "" {
		var f *os.File
		f, err = os.Open(message.DiscordFile.FilePath)

		if err != nil {
			return errors.Errorf("open png file failed, err: %s", err)
		}

		b.AddFile(message.DiscordFile.Title, message.DiscordFile.Desc, f)
	}

	if message.DiscordRobotThread != "" {
		threadID, _ := snowflake.Parse(message.DiscordRobotThread)
		_, err = client.CreateMessageInThread(b.Build(), threadID)
	} else {
		_, err = client.CreateMessage(b.Build())
	}

	if err != nil {
		return errors.Errorf("webhook.CreateMessage failed, err: %v", err)
	}

	return nil
}
