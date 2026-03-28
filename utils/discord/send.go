package discord

import (
	"context"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/pkg/errors"
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

	eb := &discord.Embed{}

	if message.AuthorName != "" {
		eb.Author = &discord.EmbedAuthor{Name: message.AuthorName}
	}

	eb.Title = message.Title

	if message.Url != "" {
		eb.URL = message.Url
	}

	eb.Description = message.Message
	eb.Timestamp = &time.Time{}
	*eb.Timestamp = time.Now()
	eb.Color = rand.Intn(0xffffff + 1)

	if message.FieldMap != nil {
		for k, v := range message.FieldMap {
			eb.Fields = append(eb.Fields, discord.EmbedField{Name: k, Value: v})
		}
	}

	b := &discord.WebhookMessageCreate{}

	//存在换行符没有生效
	if strings.Contains(message.Content, "\\n") {
		message.Content = strings.ReplaceAll(message.Content, "\\n", "\n")
	}

	b.Content = message.Content

	if message.UserName != "" {
		b.Username = message.UserName
	}

	b.Embeds = []discord.Embed{*eb}

	if message.DiscordFile != nil && message.DiscordFile.FilePath != "" {
		var f *os.File
		f, err = os.Open(message.DiscordFile.FilePath)
		if err != nil {
			return errors.Errorf("open png file failed, err: %s", err)
		}
		b.Files = []*discord.File{{Name: message.DiscordFile.Title, Description: message.DiscordFile.Desc, Reader: f}}
	}

	if message.DiscordRobotThread != "" {
		threadID, _ := snowflake.Parse(message.DiscordRobotThread)
		_, err = client.CreateMessageInThread(*b, threadID)
	} else {
		_, err = client.CreateMessage(*b, rest.CreateWebhookMessageParams{})
	}

	if err != nil {
		return errors.Errorf("webhook.CreateMessage failed, err: %v", err)
	}

	return nil
}
