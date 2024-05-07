package timeservice

import (
	"github.com/lgrisa/lib/config"
	"github.com/lgrisa/lib/utils/timeutil"
	"github.com/sirupsen/logrus"
	"time"
)

type TimeService struct {
	ctimeFunc func() time.Time
}

func NewTimeService() *TimeService {
	ctimeFunc := getCurrentTime

	if config.StartConfig.SwitchController.IsDebugMode {
		if debugCtime := config.StartConfig.SwitchController.SetServerTime; !timeutil.IsZero(debugCtime) {
			logrus.Info("设置服务器当前时间为:", debugCtime)
			diff := debugCtime.Sub(getCurrentTime())
			ctimeFunc = func() time.Time {
				return getCurrentTime().Add(diff)
			}
		}
	}

	return &TimeService{
		ctimeFunc: ctimeFunc,
	}
}

func (ts *TimeService) CurrentTime() time.Time {
	return ts.ctimeFunc()
}

func getCurrentTime() time.Time {
	return time.Now()
}
