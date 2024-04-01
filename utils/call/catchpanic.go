package call

import "github.com/sirupsen/logrus"

func CatchPanic(name string, f func()) {
	defer TryRecover(name)

	f()
}

func CatchLoopPanic(name string, f func()) {
	if name != "" {
		defer logrus.Infof("%s exit", name)
	}
	defer TryRecover(name)

	f()
}
