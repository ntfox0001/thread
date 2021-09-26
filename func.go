package thread

type CancelFunc func()
type ExecFunc func() Result
type ResultFunc func(Result)
type AsyncExecFunc func(Resolve)
