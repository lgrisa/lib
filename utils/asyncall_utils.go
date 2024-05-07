package utils

import (
	"github.com/alitto/pond"
	"github.com/lgrisa/lib/utils/log"
	"runtime/debug"
)

var (
	panicHandler = func(p interface{}) {
		stack := string(debug.Stack())
		log.LogErrorf("pool recovered from panic!!! SERIOUS PROBLEM %v %s", p, stack)
	}

	workerPool = pond.New(500, 2000, pond.Strategy(pond.Lazy()), pond.PanicHandler(panicHandler))
)

func AsyncCall(f func()) {
	if ok := workerPool.TrySubmit(f); !ok {
		f()
	}
}
