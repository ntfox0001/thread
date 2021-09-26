package thread

import (
	"context"
	"errors"
	"time"
)

type Thread struct {
	funcCase chan func()
	quit     bool
	timeout  time.Duration
}

func NewThread(funcSize uint32, timeout time.Duration) *Thread {
	return &Thread{
		funcCase: make(chan func(), funcSize),
		quit:     false,
		timeout:  timeout,
	}
}

func (t *Thread) Run() {
	go func() {
		for {
			f := <-t.funcCase
			f()

			if t.quit {
				break
			}
		}
	}()
}

func (t *Thread) RunByCatch(errFunc func(err error)) {
	defer func() {
		if err := recover(); err != nil {
			errFunc(err.(error))
			return
		}
	}()

	t.Run()
}

func (t *Thread) Close() error {
	rt := t.Exec(func() Result {
		t.quit = true
		return NewResultRt(nil)
	})

	return rt.Err
}

func (t *Thread) Then(p Promise, rtFunc ResultFunc) {
	p.Then(func(result Result) {
		err := t.AsyncExec(func() {
			rtFunc(result)
		})

		if err != nil {
			rtFunc(NewResultErr(err))
		}
	})
}

func (t *Thread) PromiseThen(p Promise, rtFunc func(Result, Resolve)) Promise {
	parent := newPromise(t.timeout)

	p.Then(func(result Result) {
		err := t.AsyncExec(func() {
			rtFunc(result, parent.resolve)
		})

		if err != nil {
			rtFunc(NewResultErr(err), parent.resolve)
		}
	})

	return parent
}

func (t *Thread) AsyncExec(f func()) error {
	ctx, cancel := t.ctx()

	select {
	case t.funcCase <- f:
		cancel()
	case <-ctx.Done():
		cancel()
		return errors.New("selectLoop-exec-timeout")
	}

	return nil
}

func (t *Thread) AsyncExecRt(f AsyncExecFunc) Promise {
	p := newPromise(t.timeout)
	err := t.AsyncExec(func() {
		f(p.resolve)
	})

	if err != nil {
		return Err(err)
	}

	return p
}

func (t *Thread) Exec(f ExecFunc) Result {
	rtChan := newChan()
	rtFunc := func() {
		rt := f()
		rtChan <- rt
	}

	err := t.AsyncExec(rtFunc)
	if err != nil {
		return NewResultErr(err)
	}

	return <-rtChan
}

func (t *Thread) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), t.timeout)
}
