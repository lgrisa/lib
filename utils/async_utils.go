package utils

import "github.com/lgrisa/lib/utils/call"

type Inviter[V any] struct {
	initChan chan struct{}
	v        V
}

func (f *Inviter[V]) Get() V {
	<-f.initChan
	return f.v
}

func StartInit[V any](name string, f func() V) *Inviter[V] {
	result := &Inviter[V]{
		initChan: make(chan struct{}),
	}
	go func() {
		defer close(result.initChan)

		result.v = f()
	}()
	return result
}

func StartInitCatchPanic[V any](name string, f func() V) *Inviter[V] {
	result := &Inviter[V]{
		initChan: make(chan struct{}),
	}
	go call.CatchPanic(name, func() {
		defer close(result.initChan)

		result.v = f()
	})
	return result
}
