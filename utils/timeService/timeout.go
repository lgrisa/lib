package ctxfunc

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Func func(ctx context.Context) (err error)

func Timeout(d time.Duration, f Func) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	return f(ctx)
}

func FuncTimeout(name string, timeout time.Duration, f func(ctx context.Context) (err error)) {
	if err := Timeout(timeout, func(ctx context.Context) (err error) {
		return f(ctx)
	}); err != nil {
		logrus.WithError(err).Error(name)
	}
}

func Timeout100ms(f Func) error {
	return Timeout(100*time.Millisecond, f)
}

func Timeout1s(f Func) error {
	return Timeout(1*time.Second, f)
}

func Timeout2s(f Func) error {
	return Timeout(2*time.Second, f)
}

func Timeout3s(f Func) error {
	return Timeout(3*time.Second, f)
}

func Timeout10s(f Func) error {
	return Timeout(10*time.Second, f)
}

func Timeout30s(f Func) error {
	return Timeout(30*time.Second, f)
}

func Timeout1m(f Func) error {
	return Timeout(60*time.Second, f)
}

func Timeout5m(f Func) error {
	return Timeout(300*time.Second, f)
}

func FuncTimeout1s(name string, f func(ctx context.Context) (err error)) {
	FuncTimeout(name, time.Second, f)
}

func FuncTimeout3s(name string, f func(ctx context.Context) (err error)) {
	FuncTimeout(name, 3*time.Second, f)
}

func FuncTimeout10s(name string, f func(ctx context.Context) (err error)) {
	FuncTimeout(name, 10*time.Second, f)
}
