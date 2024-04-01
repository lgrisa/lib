package message

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/packer/utils"
	"io"
	"net/http"
	"regexp"
	"time"
)

const (
	dingRobotURL = "https://oapi.dingtalk.com/robot/send?access_token="
)

type DingRobot struct {
	Token  string
	Secret string
}

func NewDingRobotWithSecret(token string, secret string) (*DingRobot, error) {
	if token == "" || secret == "" {
		return nil, fmt.Errorf("no token or no secret")
	}
	return &DingRobot{
		token,
		secret,
	}, nil
}

func SendMessageToDingDing(msg map[string]interface{}, token string, secret string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	header := map[string]string{
		"Content-type": "application/json",
	}

	//Millisecond time stamp
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	uri := fmt.Sprintf("%s&timestamp=%d&sign=%s", dingRobotURL+token, timestamp, sign(timestamp, secret))
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := utils.SendRequest(ctx, "POST", uri, header, b)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("send msg err: %s, token: %s, msg: %s", string(body), token, b)
	}
	return nil
}

func sign(timestamp int64, secret string) string {
	data := fmt.Sprintf("%d\n%s", timestamp, secret)

	// HMAC-SHA256 sign
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	signature := h.Sum(nil)

	// Base64 encoding
	encoded := base64.StdEncoding.EncodeToString(signature)
	return encoded

}
func (robot *DingRobot) SendMarkDownMessage(title string, text string, at ...string) error {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
	}
	if len(at) == 1 {
		if at[0] == "*" { // at all
			msg["at"] = map[string]interface{}{
				"isAtAll": true,
			}
		} else { // at specific user
			re := regexp.MustCompile(`^\+*\d{10,15}$`)
			if re.MatchString(at[0]) {
				msg["at"] = map[string]interface{}{
					"atMobiles": at,
					"isAtAll":   false,
				}
			} else {
				return fmt.Errorf(`parameter error, "at" parameter must be in "*" or mobile phone number format`)
			}
		}
	} else if len(at) > 1 {
		re := regexp.MustCompile(`^\+*\d{10,15}$`)
		for _, v := range at {
			if !re.MatchString(v) {
				return fmt.Errorf(`parameter error, "at" parameter must be in "*" or mobile phone number format`)
			}
		}
		msg["at"] = map[string]interface{}{
			"atMobiles": at,
			"isAtAll":   false,
		}
	}
	return SendMessageToDingDing(msg, robot.Token, robot.Secret)
}

func (robot *DingRobot) SendTextMessage(text string, at ...string) error {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": text,
		},
	}
	if len(at) == 1 {
		if at[0] == "*" { // at all
			msg["at"] = map[string]interface{}{
				"isAtAll": true,
			}
		} else { // at specific user
			re := regexp.MustCompile(`^\+*\d{10,15}$`)
			if re.MatchString(at[0]) {
				msg["at"] = map[string]interface{}{
					"atMobiles": at,
					"isAtAll":   false,
				}
			} else {
				return fmt.Errorf(`parameter error, "at" parameter must be in "*" or mobile phone number format`)
			}
		}
	} else if len(at) > 1 {
		re := regexp.MustCompile(`^\+*\d{10,15}$`)
		for _, v := range at {
			if !re.MatchString(v) {
				return fmt.Errorf(`parameter error, "at" parameter must be in "*" or mobile phone number format`)
			}
		}
		msg["at"] = map[string]interface{}{
			"atMobiles": at,
			"isAtAll":   false,
		}
	}
	return SendMessageToDingDing(msg, robot.Token, robot.Secret)
}
