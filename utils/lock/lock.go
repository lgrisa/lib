package lock

import (
	"context"
	"fmt"
	"github.com/lgrisa/lib/utils/concurrentmap"
	"github.com/lgrisa/lib/utils/pool"
	"github.com/lgrisa/lib/utils/reporter"
	ctxfunc "github.com/lgrisa/lib/utils/timeService"
	"github.com/lgrisa/lib/utils/timeutil"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

const (
	tickSaveInterval = 1 * time.Minute // 定时保存检查间隔
)

var (
	ErrInvalidated = errors.New("entry invalidated")
	ErrEmpty       = errors.New("entry empty")
	ErrCreateExist = errors.New("create exist entry")
	errPanic       = errors.New("provider panic")
)

type LockService[O LockObject] struct {
	provider   LockProvider[O]
	entries    *concurrentmap.Map[int64, *lockEntry[O]]
	shouldQuit chan struct{}

	saveDuration             time.Duration // 定时保存间隔
	evictNotAccessedDuration time.Duration // 定时保存间隔
}

func NewLockService[O LockObject](provider LockProvider[O], saveDuration, evictNotAccessedDuration time.Duration) *LockService[O] {
	result := &LockService[O]{
		provider:                 provider,
		entries:                  concurrentmap.New[int64, *lockEntry[O]](),
		shouldQuit:               make(chan struct{}),
		saveDuration:             saveDuration,
		evictNotAccessedDuration: evictNotAccessedDuration,
	}

	go pool.CatchLoopPanic(fmt.Sprintf("LockService.saveLoop(%v)", provider.Name()), result.saveLoop)
	go pool.CatchLoopPanic(fmt.Sprintf("LockService.evictLoop(%v)", provider.Name()), result.evictLoop)
	go pool.CatchLoopPanic(fmt.Sprintf("LockService.checkLoop(%v)", provider.Name()), result.checkLoop)
	return result
}

// 提供加载及保存数据的方法
type LockProvider[O LockObject] interface {
	Name() string                                // Service的名字
	GetObject(context.Context, int64) (O, error) // 获取某id的数据, 一般从数据库中读取并解析
	SaveObject(context.Context, int64, O) error  // 保存某id的数据, 一般调用LockObject.Marshal方法, 并把数据保存在数据库中
	CreateObject(context.Context, O) error       // 存入db中
}

// 被LockService保护的数据
type LockObject interface {
	Marshal() (pool.Buffer, error)
}

// 解锁的方法
type Unlocker func()

// 所有玩家都已确保下线, 并且所有会用到这里东西的功能模块都已正常退出之后再调用
func (s *LockService[O]) Close() {
	close(s.shouldQuit)

	logrus.Infof("LockService.%s关闭保存中", s.getProviderName())

	// 关闭时候，尝试保存4次
	for i := 0; i < 4; i++ {
		s.doEvictLoop(0, 60*time.Second) // evict everything
		if count := s.entries.Count(); count != 0 {
			logrus.WithField("count", count).Errorf("LockService.%s关闭删除了之后里面竟然还有对象... 应该其他会用到LockService的其他功能模块都完全退出之后再调用Close", s.getProviderName())
			time.Sleep(time.Second) // sleep等待下
		} else {
			break
		}
	}

	logrus.Infof("LockService.%s已关闭", s.getProviderName())
}

// 必须在hold住对象的lock时调用, 把对象标记为已删除
func (s *LockService[O]) Invalidate(id int64) {
	entry, has := s.entries.Get(id)
	if !has {
		logrus.Errorf("LockServer.%s竟然要invalidate一个不存在的entry. 必须要在Lock住对象的时候调用啊", s.getProviderName())
		return
	}

	entry.invalidated = true
}

func (s *LockService[O]) Create(id int64, data O) (Unlocker, error) {
	return s.create(id, data, true)
}

func (s *LockService[O]) Put(id int64, data O) (Unlocker, error) {
	return s.create(id, data, false)
}

func (s *LockService[O]) create(id int64, data O, insert bool) (Unlocker, error) {
	entry, has := s.entries.Get(id)
	if has {
		return nil, ErrCreateExist
	}
	// 先创建lockEntry, lock, 加入map, 在db中创建, 再unlock
	entry = &lockEntry[O]{}
	entry.Lock()
	_, ok := s.entries.SetIfAbsent(id, entry)
	if !ok {
		// 已经存在旧的, 这么巧?
		entry.Unlock() // 废弃
		return nil, ErrCreateExist
	}

	// put ok. insert into db
	if insert {
		// load 一次，看下有没有这个数据
		if _, err := s.loadObject(id); err != nil {
			if err != ErrEmpty {
				if old, removed := s.entries.Pop(id); !removed || old != entry {
					logrus.Errorf("LockService.%s.Create要删除的竟然不是原来放进去的", s.getProviderName())
				}
				entry.destroyed = true
				entry.Unlock()
				return nil, err
			}
		} else {
			if old, removed := s.entries.Pop(id); !removed || old != entry {
				logrus.Errorf("LockService.%s.Create要删除的竟然不是原来放进去的", s.getProviderName())
			}
			entry.destroyed = true
			entry.Unlock()
			return nil, ErrCreateExist
		}

		err := s.createObject(data)
		if err != nil {
			if old, removed := s.entries.Pop(id); !removed || old != entry {
				logrus.Errorf("LockService.%s.Create要删除的竟然不是原来放进去的", s.getProviderName())
			}
			entry.destroyed = true
			entry.Unlock()

			return nil, errors.Wrapf(err, "lockService.%s.createObject(%d)失败", s.getProviderName(), id)
		}
	}

	entry.data = data
	now := time.Now()
	entry.lastAccessTime, entry.lastSaveTime = now, now
	return entry.Unlock, nil
}

// 检查是否有对象在内存中
func (s *LockService[O]) IsCaching(id int64) bool {
	_, has := s.entries.Get(id)
	return has
}

// LockService关闭后还可以调用, 并没有限制. 但是修改的内容并不会保存到数据库中
// 用完之后一定要调用LockData.Unlock(), 并且在那之后不要再访问里面的内容(read也不行)
// 不unlock的话所有线程统统卡死, 也不会保存, 也不能关闭, 所有请求同一个对象的线程也卡死
func (s *LockService[O]) Lock(id int64) (o O, u Unlocker, resultErr error) {
	entry, has := s.entries.Get(id)
	if has {
		entry.Lock()
		if entry.destroyed {
			entry.Unlock()
			return s.Lock(id) // recursive
		}

		if entry.invalidated {
			entry.Unlock()
			resultErr = ErrInvalidated
			return
		}
		entry.lastAccessTime = time.Now()
		// 使用
		return entry.data, entry.Unlock, nil
	}
	// 先创建lockEntry, lock, 加入map, 去db取, 使用, 再unlock
	entry = &lockEntry[O]{}
	entry.Lock()
	old, ok := s.entries.SetIfAbsent(id, entry)
	if ok {
		// put ok. get from database and return
		data, err := s.loadObject(id)
		if err != nil {
			if old, removed := s.entries.Pop(id); !removed || old != entry {
				logrus.Errorf("LockService.%s.Lock要删除的竟然不是原来放进去的", s.getProviderName())
			}
			entry.destroyed = true
			entry.Unlock()

			if err != ErrEmpty {
				err = errors.Wrapf(err, "lockService.%s.GetObject(%d)失败", s.getProviderName(), id)
			}
			resultErr = err
			return
		}

		entry.data = data

		now := time.Now()
		entry.lastAccessTime, entry.lastSaveTime = now, now
		return data, entry.Unlock, nil
	}
	// 已经存在旧的, 这么巧?
	entry.Unlock() // 废弃
	old.Lock()
	if old.destroyed {
		// 这么巧? 刚拿出来的时候没有, 要放进去的时候已经有了, 再lock又已经destroy了?
		old.Unlock()
		return s.Lock(id) // recursive
	}
	if old.invalidated {
		old.Unlock()
		resultErr = ErrInvalidated
		return
	}
	old.lastAccessTime = time.Now()
	return old.data, old.Unlock, nil
}

// 定时删除很久没有人访问的
func (s *LockService[O]) evictLoop() {
	evictDuration := timeutil.MaxDuration(s.evictNotAccessedDuration, 10*time.Second)

	select {
	case <-time.After(1 * time.Minute):
		// 和保存loop错开运行
	case <-s.shouldQuit:
		return
	}

	// 最多2分钟
	tick := time.Tick(timeutil.MinDuration(evictDuration/2, 2*time.Minute))

	for {
		select {
		case <-s.shouldQuit:
			return

		case <-tick:
			s.doEvictLoop(evictDuration, 3*time.Second)
		}
	}
}

func (s *LockService[O]) doEvictLoop(evictTime, saveTimeout time.Duration) {
	defer pool.TryRecover(fmt.Sprintf("LockService.%s.doEvictLoop", s.getProviderName()))

	logrus.Debugf("LockService.%s.evictLoop扫描中", s.getProviderName())
	entries := s.entries.Iter()
	for ch := range entries {
		for _, item := range ch {
			entry := item.Val
			// 数据已经被加锁，跳过，下次再说
			if entry.lockTime.Load() != 0 {
				continue
			}

			entry.Lock()
			if entry.destroyed {
				entry.Unlock()
				continue
			}

			if entry.invalidated {
				if old, removed := s.entries.Pop(item.Key); !removed || old != entry {
					logrus.Errorf("LockService.%s.remove要删除的竟然不是在map中的", s.getProviderName())
				} else {
					// 已删除
					entry.destroyed = true
				}
			} else if time.Since(entry.lastAccessTime) >= evictTime || evictTime == 0 {
				// 删除时, 先lock, 检查是否destroy, 保存, 设置已destroy, unlock
				if err := s.saveObject(item.Key, entry.data, saveTimeout, true); err != nil {
					logrus.WithError(err).WithField("id", item.Key).Errorf("LockService.%s要删除LockData, 保存时出错", s.getProviderName())
				} else {
					// 已保存
					if old, removed := s.entries.Pop(item.Key); !removed || old != entry {
						logrus.Errorf("LockService.%s.remove要删除的竟然不是在map中的", s.getProviderName())
					} else {
						// 已删除
						entry.destroyed = true
					}
				}
			}
			entry.Unlock()
		}
	}
	logrus.Debugf("LockService.%s.evictLoop扫描完成", s.getProviderName())
}

// 定时保存
func (s *LockService[O]) saveLoop() {
	tick := time.Tick(tickSaveInterval)

	for {
		select {
		case <-s.shouldQuit:
			return

		case <-tick:
			s.doSaveLoop(s.saveDuration)
		}
	}
}

func (s *LockService[O]) doSaveLoop(saveInterval time.Duration) {
	defer pool.TryRecover(fmt.Sprintf("LockService.%s.doSaveLoop", s.getProviderName()))
	// 不直接取这里的时间. 可能db很慢很慢, 一个循环花很久, 导致不停在保存, 更增加了db压力
	// 全量扫描保存
	logrus.Debugf("LockService.%s保存检查中", s.getProviderName())
	entries := s.entries.Iter()
	for ch := range entries {
		for _, item := range ch {
			entry := item.Val
			// 数据已经被加锁，跳过，下次再说
			if entry.lockTime.Load() != 0 {
				continue
			}

			entry.Lock()
			if entry.destroyed || entry.invalidated {
				entry.Unlock()
				continue
			}
			if entry.lastAccessTime.After(entry.lastSaveTime.Add(-time.Second)) && time.Since(entry.lastSaveTime) >= saveInterval {
				// 必须是上次被人接触过的时间 > (上次保存的时间 - 1秒) && 超出保存间隔
				// 上次保存之后的1秒内有用过, 也重新保存
				if err := s.saveObject(item.Key, entry.data, 3*time.Second, true); err != nil {
					logrus.WithError(err).WithField("id", item.Key).Errorf("LockService.%s定时保存出错", s.getProviderName())
				} else {
					// 已保存
					entry.lastSaveTime = time.Now()
				}
			}
			entry.Unlock()
		}
	}
	logrus.Debugf("LockService.%s保存检查完成", s.getProviderName())
}

func (s *LockService[O]) getProviderName() string {
	return s.provider.Name()
}

func (s *LockService[O]) createObject(o O) (err error) {
	defer pool.RecoverFunc(fmt.Sprintf("LockService.%s.createObject", s.getProviderName()), func() {
		err = errPanic
	})

	return ctxfunc.Timeout10s(func(ctx context.Context) (err error) {
		return s.provider.CreateObject(ctx, o)
	})
}

func (s *LockService[O]) loadObject(key int64) (o O, err error) {
	defer pool.RecoverFunc(fmt.Sprintf("LockService.%s.loadObject", s.getProviderName()), func() {
		err = errPanic
	})

	err = ctxfunc.Timeout10s(func(ctx context.Context) (err error) {
		o, err = s.provider.GetObject(ctx, key)
		return err
	})
	return
}

func (s *LockService[O]) saveObject(key int64, data O, timeout time.Duration, sort bool) (err error) {
	defer pool.RecoverFunc(fmt.Sprintf("LockService.%s.saveObject", s.getProviderName()), func() {
		err = errPanic
	})

	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return ctxfunc.Timeout(timeout, func(ctx context.Context) (err error) {
		return s.provider.SaveObject(ctx, key, data)
	})
}

// 检查所有的entry是否有长时间没有unlock的
func (s *LockService[O]) checkLoop() {
	<-time.After(30 * time.Second) // 和saveLoop & evictLoop 错开
	tick := time.Tick(2 * time.Minute)

	for {
		select {
		case <-s.shouldQuit:
			return

		case <-tick:
			s.doCheckLoop()
		}
	}
}

func (s *LockService[O]) doCheckLoop() {
	shouldQuit := make(chan struct{})

	var checkEntry *lockEntry[O]
	entries := s.entries.Iter()
	timeout := time.After(5 * time.Second)

	go pool.CatchPanic("LockService.doCheckLoop", func() {
		select {
		case <-shouldQuit:
			return
		case <-timeout:
		}
		if checkEntry := checkEntry; checkEntry != nil {
			for i := 0; i < 100; i++ {
				logrus.WithField("object", checkEntry).Errorf("LockService.%s 检测到长时间没有unlock的对象!!!! SEVERE!!!", s.getProviderName())
			}
			dumpStacks(s.getProviderName())
		}
	})
	for ch := range entries {
		for _, item := range ch {
			checkEntry = item.Val
			checkEntry.Lock()
			checkEntry.Unlock()
		}
	}
	close(shouldQuit)
}

func dumpStacks(name string) {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, true)]
	ioutil.WriteFile(name+"_"+time.Now().Format(timeutil.SecondsLayout), buf, os.ModePerm)

	reporter.FormatStack(string(buf), "服务器死锁检查")
}

type lockEntry[O LockObject] struct {
	mux            sync.Mutex
	lockTime       atomic.Int64
	lastAccessTime time.Time // 上次有人请求的时间
	lastSaveTime   time.Time // 上次自动保存的时间
	destroyed      bool      // 是否已废弃失效
	invalidated    bool      // 是否已被标记为从db中删除
	data           O         // 真正的内容
}

func (l *lockEntry[O]) Lock() {
	l.mux.Lock()
	l.lockTime.Store(time.Now().Unix())
}

func (l *lockEntry[O]) Unlock() {
	l.lockTime.Store(0)
	l.mux.Unlock()
}

func (l lockEntry[O]) String() string {
	return fmt.Sprintf("%v", l.data)
}
