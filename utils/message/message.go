package message

import (
	"github.com/lgrisa/library/utils"
	"github.com/lgrisa/library/utils/call"
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
	utils.LogErrorf("SendMessage: %v", msg)

	go call.CatchLoopPanic("SendMessage", func() {
		message := m.messagePrefix + msg

		if m.DingRobot != nil {
			err := m.DingRobot.SendTextMessage(message)
			if err != nil {
				utils.LogErrorf("SendTextMessage DingDing error: %v", err)
			}
		}

		if m.discordWebhookUrl != "" {
			sendDiscord(m.discordWebhookUrl, m.discordRobotThread, m.discordMessageTitle, message)
		}
	})
}
