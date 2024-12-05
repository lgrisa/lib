package dingding

const (
	dingRobotURL = "https://oapi.dingtalk.com/robot/send?access_token="
)

type DingRobot struct {
	Token  string
	Secret string
}
