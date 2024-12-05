package message

import (
	"github.com/lgrisa/lib/utils"
	"github.com/lgrisa/lib/utils/call"
	"github.com/lgrisa/lib/utils/dingding"
	"github.com/lgrisa/lib/utils/discord"
	"github.com/pkg/errors"
)

type SendMessageClient struct {
	messagePrefix       string
	discordWebhookUrl   string
	discordRobotThread  string
	discordMessageTitle string
	DingRobot           *dingding.DingRobot
}

func NewSendMessageClient(messagePrefix, dingDingAccessToken, dingDingSecret, discordWebhookUrl string) (*SendMessageClient, error) {
	sendClient := &SendMessageClient{}

	if dingDingAccessToken != "" && dingDingSecret != "" {
		dingTalkCli, err := dingding.NewDingRobotWithSecret(dingDingAccessToken, dingDingSecret)
		if err != nil {
			return nil, errors.Errorf("NewDingRobotWithSecret error: %v", err)
		}

		sendClient.DingRobot = dingTalkCli
	}

	return sendClient, nil
}

func (m *SendMessageClient) SendTextMessage(msg string) {
	utils.LogPrintf("SendMessage: %v", msg)

	go call.CatchLoopPanic("SendMessage", func() {
		message := m.messagePrefix + msg

		if m.DingRobot != nil {
			err := m.DingRobot.SendTextMessage(message)
			if err != nil {
				utils.LogErrorF("SendTextMessage DingDing error: %v", err)
			}
		}

		if m.discordWebhookUrl != "" {
			if err := discord.SendDiscordNoGoroutines(
				&discord.Message{
					DiscordWebhookUrl:  m.discordWebhookUrl,
					DiscordRobotThread: m.discordRobotThread,
					Title:              m.discordMessageTitle,
					Message:            message,
				}); err != nil {
				utils.LogErrorF("SendDiscord error: %v", err)
			}
		}
	})
}
