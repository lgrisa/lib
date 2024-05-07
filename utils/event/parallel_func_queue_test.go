package event

import (
	. "github.com/onsi/gomega"
	"sync/atomic"
	"testing"
	"time"
)

type Int32 struct {
	v int32
}

func (i *Int32) Load() int32 {
	return atomic.LoadInt32(&i.v)
}

func (i *Int32) Inc() int32 {
	return atomic.AddInt32(&i.v, 1)
}

func TestParallelFuncQueue(t *testing.T) {
	RegisterTestingT(t)

	queue := NewParallelFuncQueue(10, 5, "test")

	counter := Int32{}
	for i := 0; i < 1000; i++ {
		queue.MustFunc(func() {
			counter.Inc()
		})
	}

	Eventually(func() int32 {
		return counter.Load()
	}).WithTimeout(10 * time.Second).
		Should(Equal(int32(1000)))

	for i := 0; i < 1000; i++ {
		queue.MustFunc(func() {
			counter.Inc()
		})
	}

	queue.Close(false)

	Ω(counter.Load()).Should(BeEquivalentTo(2000))

}
