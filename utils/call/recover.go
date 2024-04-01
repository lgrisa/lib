package call

import (
	"fmt"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

var server string

func SetServer(s string) {
	server = s
}

func GetServer() string {
	return server
}

func TryRecover(handlerName string) bool {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		logrus.WithField("err", r).Error(handlerName + " recovered from panic!!! SERIOUS PROBLEM " + stack)
		fmt.Println(r, stack)
		return true
	}
	return false
}

func RecoverFunc(handlerName string, f func()) bool {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		logrus.WithField("err", r).WithField("stack", stack).Error(handlerName + " recovered from panic!!! SERIOUS PROBLEM")
		fmt.Println(r, stack)
		f()
		return true
	}
	return false
}

func RecoverError(handlerName string, f func(err interface{})) bool {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		logrus.WithField("err", r).WithField("stack", stack).Error(handlerName + " recovered from panic!!! SERIOUS PROBLEM")
		fmt.Println(r, stack)
		f(r)
		return true
	}
	return false
}

func TryRecoverWithFunc(name string, f func(hasRecover bool)) {
	hasRecover := TryRecover(name)
	f(hasRecover)
}
