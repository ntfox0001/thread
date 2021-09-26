package thread

import (
	"context"
	"sync"
	"time"
)

type Ticker struct {
	selectLoop IThread
	cancels    []context.CancelFunc
	waitGroup  sync.WaitGroup
	logger     ILogger
}

func NewTicker(selectLoop IThread, logger ILogger) *Ticker {
	return &Ticker{
		selectLoop: selectLoop,
		cancels:    make([]context.CancelFunc, 0),
		waitGroup:  sync.WaitGroup{},
		logger:     logger,
	}
}

func (t *Ticker) Start(interval time.Duration, f func()) {
	ctx, cancel := context.WithCancel(context.Background())

	go t.run(ctx, interval, f)
	t.cancels = append(t.cancels, cancel)
}

func (t *Ticker) Release() Promise {
	for _, c := range t.cancels {
		c()
	}

	t.cancels = make([]context.CancelFunc, 0)

	p := NewPromise(func(resolve Resolve) {
		t.waitGroup.Wait()
		resolve.R(nil)
	})

	return p
}

func (t *Ticker) run(ctx context.Context, interval time.Duration, f func()) {
	t.waitGroup.Add(1)

	ticker := time.NewTicker(interval)
runnable:
	for {
		select {
		case <-ctx.Done():
			break runnable
		case <-ticker.C:
			// 在目标selectLoop上执行函数
			err := t.selectLoop.AsyncExec(f)

			if err != nil {
				t.logf("Ticker: call error: %s\n", err.Error())
			}
		}
	}

	t.waitGroup.Done()
}

func (t *Ticker) logf(format string, args ...interface{}) {
	if t.logger != nil {
		t.logger.Errorf(format, args...)
	}
}
