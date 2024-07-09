package message

import (
	"github.com/lgrisa/lib/utils"
	"github.com/lgrisa/lib/utils/call"
	"github.com/pkg/errors"
)

type SendMessageClient struct {
	messagePrefix       string
	discordWebhookUrl   string
	discordRobotThread  string
	discordMessageTitle string
	DingRobot           *DingRobot
}

func NewSendMessageClient(messagePrefix, dingDingAccessToken, dingDingSecret, discordWebhookUrl string) (*SendMessageClient, error) {
	sendClient := &SendMessageClient{}

	if dingDingAccessToken != "" && dingDingSecret != "" {
		dingTalkCli, err := NewDingRobotWithSecret(dingDingAccessToken, dingDingSecret)
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
			if err := SendDiscordNoGoroutines(
				&DiscordMessage{
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
