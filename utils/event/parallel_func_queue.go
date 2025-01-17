package event

import (
	"container/list"
	"fmt"
	"github.com/lgrisa/lib/utils/pool"
	"sync"
	"time"
)

type ParallelFuncQueue struct {
	name string

	handle    bool
	count     int
	funcQueue chan func()

	funcCache    *list.List
	funcCacheMux sync.Mutex

	workers []*worker

	queueEmptyNotify chan struct{}
	closeNotify      chan struct{}
	loopExitNotify   chan struct{}
}

func NewParallelFuncQueue(n, workerCount uint64, name string) *ParallelFuncQueue {
	q := &ParallelFuncQueue{
		name:             name,
		count:            int(n),
		funcQueue:        make(chan func(), n),
		funcCache:        list.New(),
		queueEmptyNotify: make(chan struct{}),
		closeNotify:      make(chan struct{}),
		loopExitNotify:   make(chan struct{}),
	}

	if workerCount <= 0 {
		workerCount = 1
	}

	q.workers = make([]*worker, workerCount)
	for i := uint64(0); i < workerCount; i++ {
		w := newWorker(fmt.Sprintf("%s-%d", name, i), q.funcQueue, q.onQueueEmpty)
		q.workers[i] = w
	}

	go pool.CatchLoopPanic(name, q.loop)

	return q
}

// Close immediately true表示不处理队列中剩余的内容
func (s *ParallelFuncQueue) Close(immediately bool) {
	close(s.closeNotify)
	<-s.loopExitNotify

	// 停掉worker
	for _, w := range s.workers {
		w.close()
	}

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

func (s *ParallelFuncQueue) onQueueEmpty() {
	select {
	case s.queueEmptyNotify <- struct{}{}:
	default:
	}
}

func (s *ParallelFuncQueue) loop() {
	defer close(s.loopExitNotify)

	secondTicker := time.NewTicker(time.Second)

	for {
		select {
		case <-s.queueEmptyNotify:
			// 队列中没有数据，处理缓存中的数据
			if len(s.funcQueue) <= 0 {
				s.handleCache()
			}

		case <-secondTicker.C:
			// 定时看一下队列中有没有数据
			if len(s.funcQueue) <= 0 {
				s.handleCache()
			}

		case <-s.closeNotify:
			return
		}
	}
}

func (s *ParallelFuncQueue) TryFunc(f func()) bool {
	select {
	case s.funcQueue <- f:
		return true
	case <-s.closeNotify:
		return false
	default:
		return false
	}
}

func (s *ParallelFuncQueue) TimeoutFunc(f func(), timeout time.Duration) bool {
	if s.TryFunc(f) {
		return true
	}

	select {
	case s.funcQueue <- f:
		return true
	case <-time.After(timeout):
		return false
	case <-s.closeNotify:
		return false
	}
}

func (s *ParallelFuncQueue) MustFunc(f func()) {
	if !s.TryFunc(f) {
		s.pushToCache(f)
	}
}

func (s *ParallelFuncQueue) pushToCache(f func()) {
	s.funcCacheMux.Lock()
	defer s.funcCacheMux.Unlock()
	s.funcCache.PushBack(f)
}

func (s *ParallelFuncQueue) popCache(n int) []func() {
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

func (s *ParallelFuncQueue) handleCache() {
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

func newWorker(name string, funcQueue chan func(), onFuncQueueEmpty func()) *worker {
	w := &worker{
		name:             name,
		funcQueue:        funcQueue,
		onFuncQueueEmpty: onFuncQueueEmpty,
		closeNotify:      make(chan struct{}),
		loopExitNotify:   make(chan struct{}),
	}

	go pool.CatchLoopPanic(name, w.loop)

	return w
}

// 每个worker一个线程，负责处理func
type worker struct {
	name      string
	funcQueue chan func()

	onFuncQueueEmpty func()

	closeNotify    chan struct{}
	loopExitNotify chan struct{}
}

func (s *worker) close() {
	close(s.closeNotify)
	<-s.loopExitNotify
}

func (s *worker) loop() {
	defer close(s.loopExitNotify)

	for {
		select {
		case f := <-s.funcQueue:
			pool.CatchPanic(s.name, f)

			if len(s.funcQueue) <= 0 {
				s.onFuncQueueEmpty()
			}

		case <-s.closeNotify:
			return
		}
	}
}
