package call

import (
	"github.com/lgrisa/library/utils"
)

func CatchPanic(name string, f func()) {
	defer TryRecover(name)

	f()
}

func CatchLoopPanic(name string, f func()) {
	if name != "" {
		defer utils.LogInfof("%s exit", name)
	}
	defer TryRecover(name)

	f()
}
