package call

import (
	"fmt"
	"testing"
)

func TestRecover(t *testing.T) {
	fmt.Println(1)
	tryRecover()
	fmt.Println(2)
	recoverFunc()
	fmt.Println(3)
	recoverError()
	fmt.Println(4)

	CatchPanic("CP", doPanic)
	fmt.Println(6)
	CatchLoopPanic("CLP",doPanic )
	fmt.Println(7)
}

func doPanic() {
	var array []int
	fmt.Println(array[1])
}

func tryRecover() {
	defer TryRecover("tryRecover")
	doPanic()
}

func recoverFunc() {
	defer RecoverFunc("tryRecover", func() {
		fmt.Println(5)
	})
	doPanic()
}

func recoverError() {
	defer RecoverError("tryRecover", func(err interface{}) {
		fmt.Println(err)
	})
	doPanic()
}
