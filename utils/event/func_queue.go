package event

import (
	"container/list"
	"github.com/lgrisa/lib/utils/pool"
	"sync"
	"time"
)

type FuncQueue struct {
	name string

	handle    bool
	count     int
	funcQueue chan func()

	funcCache    *list.List
	funcCacheMux sync.Mutex

	closeNotify    chan struct{}
	loopExitNotify chan struct{}
}

func NewFuncQueue(n uint64, name string) *FuncQueue {
	q := &FuncQueue{
		name:           name,
		count:          int(n),
		funcQueue:      make(chan func(), n),
		funcCache:      list.New(),
		closeNotify:    make(chan struct{}),
		loopExitNotify: make(chan struct{}),
	}

	go pool.CatchLoopPanic(name, q.loop)

	return q
}

func (s *FuncQueue) loop() {
	defer close(s.loopExitNotify)

	for {
		select {
		case f := <-s.funcQueue:
			pool.CatchPanic(s.name, f)

			if len(s.funcQueue) <= 0 {
				// 队列中没有数据，处理缓存中的数据
				s.handleCache()
			}

		case <-s.closeNotify:
			return
		}
	}
}

// Close immediately true表示不处理队列中剩余的内容
func (s *FuncQueue) Close(immediately bool) {
	close(s.closeNotify)
	<-s.loopExitNotify

	if !immediately {
		// 关闭时候，需要处理队列剩余的内容
	out:
		for {
			select {
			case f := <-s.funcQueue:
				pool.CatchPanic(s.name, f)
			default:
				break out
			}
		}

		s.funcCacheMux.Lock()
		defer s.funcCacheMux.Unlock()

		for e := s.funcCache.Front(); e != nil; {
			next := e.Next()
			value := s.funcCache.Remove(e)
			e = next

			if f, ok := value.(func()); ok {
				pool.CatchPanic(s.name, f)
			}
		}
	}

}

func (s *FuncQueue) TryFunc(f func()) bool {
	select {
	case s.funcQueue <- f:
		return true
	default:
		return false
	}
}

func (s *FuncQueue) TimeoutFunc(f func(), timeout time.Duration) bool {
	if s.TryFunc(f) {
		return true
	}

	select {
	case s.funcQueue <- f:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (s *FuncQueue) MustFunc(f func()) {
	if !s.TryFunc(f) {
		s.pushToCache(f)
	}
}

func (s *FuncQueue) pushToCache(f func()) {
	s.funcCacheMux.Lock()
	defer s.funcCacheMux.Unlock()
	s.funcCache.PushBack(f)
}

func (s *FuncQueue) popCache(n int) []func() {
	s.funcCacheMux.Lock()
	defer s.funcCacheMux.Unlock()

	var fs []func()
	count := 0
	for e := s.funcCache.Front(); e != nil; {
		next := e.Next()
		value := s.funcCache.Remove(e)
		e = next

		if f, ok := value.(func()); ok {
			fs = append(fs, f)

			count++
			if count >= n {
				break
			}
		}
	}

	return fs
}

func (s *FuncQueue) handleCache() {
	// 每次添加到80%
	n := s.count * 8 / 10
	if n < 10 {
		n = 10
	}

	fs := s.popCache(n)

	for i, f := range fs {
		if !s.TryFunc(f) {
			// 添加不进去，加回到列表中，从尾巴开始加
			fsn := len(fs) - i
			func() {
				s.funcCacheMux.Lock()
				defer s.funcCacheMux.Unlock()

				for i := 0; i < fsn; i++ {
					idx := fsn - i - 1
					f := fs[idx]
					s.funcCache.PushFront(f)
				}
			}()
			break
		}
	}
}
