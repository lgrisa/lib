package reporter

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var defaultFeishuReporter *feishuReporter

func InitFeishu() {
	token := startconfig.StartConfig.ErrReporterToken
	if token == "" {
		return
	}

	defaultFeishuReporter = &feishuReporter{
		name:  hostname.GetValue(),
		token: token,
	}
}

func Format(format string, args ...interface{}) {
	Text(fmt.Sprintf(format, args...))
}

func FormatStack(stack, format string, args ...interface{}) {
	TextStack(fmt.Sprintf(format, args...), stack)
}

func Text(text string) {
	if defaultFeishuReporter == nil {
		return
	}

	stack := debug.Stack()
	if len(stack) > 1000 {
		stack = stack[:1000]
	}
	stackString := string(stack)

	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			logrus.WithField("err", r).Error("reporter recovered from panic!!! SERIOUS PROBLEM " + stack)
			fmt.Println(r, stack)
		}
	}()
	defaultFeishuReporter.ReportText(text, stackString)
}

func TextStack(text, stackString string) {
	if defaultFeishuReporter == nil {
		return
	}

	if len(stackString) > 1000 {
		stackString = stackString[:1000]
	}

	if r := recover(); r != nil {
		stack := string(debug.Stack())
		logrus.WithField("err", r).Error("reporter recovered from panic!!! SERIOUS PROBLEM " + stack)
		fmt.Println(r, stack)
	}
	defaultFeishuReporter.ReportText(text, stackString)

}

type feishuReporter struct {
	name  string
	token string
}

type message struct {
	MsgType string `json:"msg_type"`

	Content content `json:"content"`
}

type content struct {
	Text string `json:"text"`
}

func (r *feishuReporter) ReportText(text, stack string) {
	sb := strings.Builder{}
	sb.WriteString("uestarserver report, server: ")
	sb.WriteString(r.name)
	sb.WriteString("\n")
	sb.WriteString(text)
	sb.WriteString("\n")
	sb.WriteString(stack)

	msg := &message{}
	msg.MsgType = "text"
	msg.Content.Text = sb.String()

	postData, err := jsoniter.MarshalToString(msg)
	if err != nil {
		logrus.Errorf("report to feishu failed, json marshal error: %v", err)
		return
	}

	_, err = http.Post(r.token, "application/json", strings.NewReader(postData))
	if err != nil {
		logrus.WithError(err).Debugf("report to feishu failed")
		return
	}
}
