package thread

import "time"

type IThread interface {
	AsyncExec(f func()) error
	Exec(f ExecFunc) Result
	AsyncExecRt(f AsyncExecFunc) Promise
	Then(p Promise, rtFunc ResultFunc)
	PromiseThen(p Promise, rtFunc func(Result, Resolve)) Promise
	Close() error
}

// ITimer 提供线程内timer
type ITimer interface {
	AddTimer(duration time.Duration, f func()) (CancelFunc, error)
}

type ILogger interface {
	Errorf(format string, args ...interface{})
}
