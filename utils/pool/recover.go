package pool

import (
	"github.com/lgrisa/lib/utils/logutil"
	"runtime/debug"
)

func TryRecover(handlerName string) bool {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		logutil.LogErrorF(handlerName+" recovered from panic!!! SERIOUS PROBLEM, err: %v, stack: %v", r, stack)
		return true
	}
	return false
}

func RecoverFunc(handlerName string, f func()) bool {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		logutil.LogErrorF(handlerName+" recovered from panic!!! SERIOUS PROBLEM, err: %v, stack: %v", r, stack)
		f()
		return true
	}
	return false
}

func RecoverError(handlerName string, f func(err interface{})) bool {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		logutil.LogErrorF(handlerName+" recovered from panic!!! SERIOUS PROBLEM, err: %v, stack: %v", r, stack)
		f(r)
		return true
	}
	return false
}

func TryRecoverWithFunc(name string, f func(hasRecover bool)) {
	hasRecover := TryRecover(name)
	f(hasRecover)
}
