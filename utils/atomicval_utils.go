package utils

import (
	uatomic "go.uber.org/atomic"
	"sync/atomic"
)

func Get[T any](val *atomic.Value) *T {
	if ref := val.Load(); ref == nil {
		return nil
	} else {
		return ref.(*T)
	}
}

func GetV[T any](val *uatomic.Value) *T {
	if ref := val.Load(); ref == nil {
		return nil
	} else {
		return ref.(*T)
	}
}
