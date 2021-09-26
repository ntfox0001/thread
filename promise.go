package thread

import (
	"context"
	"errors"
	"time"
)

type Result struct {
	Val interface{}
	Err error
}

func NewResultErr(err error) Result {
	return Result{nil, err}
}
func NewResultRt(val interface{}) Result {
	return Result{val, nil}
}

// ------------------------------------------------------------------------------------------------------

type Resolve struct {
	channel chan Result
	timeout time.Duration
}

func NewResolve(timeout time.Duration) Resolve {
	return Resolve{
		channel: make(chan Result, 1),
		timeout: timeout,
	}
}

func (r Resolve) Return(et *Result) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	select {
	case r.channel <- *et:
	case <-ctx.Done():
	}

	cancel()
}

// R 返回指定的值
func (r Resolve) R(val interface{}) {
	r.Return(&Result{val, nil})
}

// E 返回错误值
func (r Resolve) E(err error) {
	r.Return(&Result{nil, err})
}

// Nil 返回空值
func (r Resolve) Nil() {
	r.Return(&Result{nil, nil})
}

// ------------------------------------------------------------------------------------------------------
type Promise struct {
	resolve Resolve
	timeout time.Duration
}

const defaultTimeout = time.Second * 10

func newPromise(timeout time.Duration) Promise {
	return Promise{
		resolve: NewResolve(timeout),
		timeout: timeout,
	}
}

// newChan 创建一个异步应答
func newChan() chan Result {
	ac := make(chan Result, 1)
	return ac
}

func Err(err error) Promise {
	p := Promise{
		resolve: NewResolve(defaultTimeout),
		timeout: defaultTimeout,
	}

	p.resolve.E(err)

	return p
}

func Rt(val interface{}) Promise {
	p := Promise{
		resolve: NewResolve(defaultTimeout),
		timeout: defaultTimeout,
	}

	p.resolve.R(val)

	return p
}

// NewPromise 开始一个异步请求
func NewPromise(f func(resolve Resolve)) Promise {
	p := Promise{
		resolve: NewResolve(defaultTimeout),
		timeout: defaultTimeout,
	}

	go func() {
		f(p.resolve)
	}()

	return p
}

// SyncThen 阻塞等待
func (p Promise) SyncThen() Result {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)

	select {
	case <-ctx.Done():
		cancel()
		return NewResultErr(errors.New("receive-timeout"))
	case rt := <-p.resolve.channel:
		cancel()
		return rt
	}
}

// Then 异步接受返回
func (p Promise) Then(rtFunc ResultFunc) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), p.timeout)

		select {
		case <-ctx.Done():
			rtFunc(NewResultErr(errors.New("receive-timeout")))
		case rt := <-p.resolve.channel:
			rtFunc(rt)
		}

		cancel()
	}()
}
