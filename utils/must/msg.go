package must

import (
	"github.com/lgrisa/lib/utils/pbutil"
	"github.com/sirupsen/logrus"
)

func Msg(data pbutil.Buffer, err error) pbutil.Buffer {
	if err != nil {
		logrus.WithError(err).Errorf("must.Msg fail")
		return pbutil.Empty
	}
	return data
}
