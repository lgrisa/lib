package dingding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lgrisa/lib/utils"
	"net/http"
	"regexp"
	"time"
)

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

	_, code, respBody, err := utils.Request(ctx, "POST", uri, header, bytes.NewReader(b))
	if code != http.StatusOK {
		return fmt.Errorf("send msg err: %s, token: %s, msg: %s", string(respBody), token, b)
	}

	return nil
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
