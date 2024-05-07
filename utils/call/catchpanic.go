package call

import (
	"github.com/lgrisa/lib/utils/log"
)

func CatchPanic(name string, f func()) {
	defer TryRecover(name)

	f()
}

func CatchLoopPanic(name string, f func()) {
	if name != "" {
		defer log.LogInfof("%s exit", name)
	}
	defer TryRecover(name)

	f()
}
