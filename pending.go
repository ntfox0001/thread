package thread

// Pending 当发起多个异步请求时，某种情况下，会连续发出同样的请求，
// Pending可以合并这些请求并在第一个回应返回时，顺序调用处理函数
// 调用者首先要自己判断请求是否重复，然后将处理函数Attach进Pending
// 如使用一个map[Type]Pending的结构，对于同一种请求，使用相同Type
type Pending struct {
	resultFuncs []ResultFunc
	promise     Promise
	final       Handler
}

type Handler interface {
	Invoke()
}

type HandlerFunc func()

func (f HandlerFunc) Invoke() {
	f()
}

func NewPending(p Promise) *Pending {
	return &Pending{
		resultFuncs: make([]ResultFunc, 0),
		promise:     p,
		final:       nil,
	}
}

func (p *Pending) Attach(rtFunc ResultFunc) {
	p.resultFuncs = append(p.resultFuncs, rtFunc)
}

// SetFinal 设置一个最终处理函数，在所有处理函数之后调用
func (p *Pending) SetFinal(f Handler) {
	p.final = f
}

func (p *Pending) Then(sl IThread) {
	sl.Then(p.promise, func(result Result) {
		for i := range p.resultFuncs {
			p.resultFuncs[i](result)
		}

		if p.final != nil {
			p.final.Invoke()
		}
	})
}
