package event

import (
	"github.com/lgrisa/lib/utils/call"
	"github.com/lgrisa/lib/utils/math/i64"
	"github.com/lgrisa/lib/utils/timer"
	"github.com/sirupsen/logrus"
	"time"
)

type EventQueue struct {
	// 事件处理队列
	funcChan         chan *eventAction
	closeNotifier    chan struct{}
	loopExitNotifier chan struct{}

	queueTimeoutWheel *timer.TimingWheel
	queueTimeout      time.Duration

	name string
}

func NewEventQueue(queueLength uint64, timeout time.Duration, name string) *EventQueue {
	if timeout <= 0 {
		logrus.Panicf("NewEventQueue timeout <= 0")
	}

	m := &EventQueue{}

	m.funcChan = make(chan *eventAction, queueLength)
	m.closeNotifier = make(chan struct{})
	m.loopExitNotifier = make(chan struct{})

	interval := 500 * time.Millisecond
	buckets := int(i64.DivideTimes(timeout.Nanoseconds(), interval.Nanoseconds()) + 1)
	m.queueTimeoutWheel = timer.NewTimingWheel(interval, buckets)
	m.queueTimeout = timeout

	go call.CatchLoopPanic(name, m.loop)

	return m
}

// TimeoutFunc 放进去直到超时
func (r *EventQueue) TimeoutFunc(waitResult bool, f func()) (funcCalled bool) {
	e := &eventAction{f: f, called: make(chan struct{})}

	select {
	case r.funcChan <- e:
		if waitResult {
			select {
			case <-r.loopExitNotifier:
				return false // main loop exit

			case <-e.called:
				return true
			}
		} else {
			return true // put success
		}

	case <-r.queueTimeoutWheel.After(r.queueTimeout):
		return false

	case <-r.closeNotifier:
		return false
	}
}

// Func 放进去位置
func (r *EventQueue) Func(waitResult bool, f func()) (funcCalled bool) {
	e := &eventAction{f: f, called: make(chan struct{})}

	select {
	case r.funcChan <- e:
		if waitResult {
			select {
			case <-r.loopExitNotifier:
				return false // main loop exit

			case <-e.called:
				return true
			}
		} else {
			return true // put success
		}

	case <-r.closeNotifier:
		return false
	}
}

func (r *EventQueue) loop() {

	defer close(r.loopExitNotifier)

	for {
		select {
		case f := <-r.funcChan:
			call.CatchPanic(r.name, f.f)
			close(f.called)
		case <-r.closeNotifier:
			return
		}
	}

}

func (r *EventQueue) Stop() {
	close(r.closeNotifier)

	<-r.loopExitNotifier
}

type eventAction struct {
	f      func()
	called chan struct{}
}
