package pool

import (
	"github.com/lgrisa/lib/utils/logutil"
	"github.com/panjf2000/ants/v2"
	"reflect"
	"sync"
)

var (
	gWorkerPool *ants.Pool
	gWorkerOnce sync.Once
	gWorkWait   sync.WaitGroup
	gWordStopCh = make(chan struct{})
)

func InitGroupWorkerPool(size int) {
	gWorkerOnce.Do(func() {
		gWorkerPool, _ = ants.NewPool(size,
			ants.WithPreAlloc(true),     //是否预分配协程池的内存
			ants.WithNonblocking(false), //如果任务无法提交，直接返回错误
			ants.WithExpiryDuration(0),
			//ants.WithPanicHandler(func(i interface{}) { //提供自定义的 Panic 处理逻辑
			//}),
			ants.WithLogger(logutil.GetLogger()),
			ants.WithMaxBlockingTasks(0),
		)

		gWorkWait.Add(1)

		Go(func() {
			defer gWorkWait.Done()
			<-gWordStopCh
			gWorkerPool.Release()
		})
	})
}

func Go(f func()) {

	// Submit 提交一个任务到协程池
	// 如果协程池已经关闭，Submit 会返回一个错误
	gWorkWait.Add(1)
	if err := gWorkerPool.Submit(func() {
		defer gWorkWait.Done()
		f()
	}); err != nil {
		gWorkWait.Done()

		logutil.LogErrorF("workerPool.Submit f name: %s, err: %v", reflect.TypeOf(f).Name(), err)
	}
}

// 附带停止信号

func GoWithStop(f func(stopCh <-chan struct{})) {
	gWorkWait.Add(1)
	if err := gWorkerPool.Submit(func() {
		defer gWorkWait.Done()
		f(gWordStopCh)
	}); err != nil {
		gWorkWait.Done()

		logutil.LogErrorF("workerPool.Submit f name: %s, err: %v", reflect.TypeOf(f).Name(), err)
	}
}
