package call

import (
	"github.com/disgoorg/log"
)

func CatchPanic(name string, f func()) {
	defer TryRecover(name)

	f()
}

func CatchLoopPanic(name string, f func()) {
	if name != "" {
		defer log.Infof("%s exit", name)
	}
	defer TryRecover(name)

	f()
}
