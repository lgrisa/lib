package utils

import (
	"github.com/alitto/pond"
	"runtime/debug"
)

var (
	panicHandler = func(p interface{}) {
		stack := string(debug.Stack())
		LogErrorF("pool recovered from panic!!! SERIOUS PROBLEM %v %s", p, stack)
	}

	workerPool = pond.New(500, 2000, pond.Strategy(pond.Lazy()), pond.PanicHandler(panicHandler))
)

func AsyncCall(f func()) {
	if ok := workerPool.TrySubmit(f); !ok {
		f()
	}
}
