package thread

import (
	"sync/atomic"
)

// All 接受多个Promise，返回一个[]interface{},如果有错误发生，那么返回第一个错误
func All(pl ...Promise) Promise {
	return NewPromise(func(resolve Resolve) {
		plSize := len(pl)
		rts := make([]interface{}, plSize)
		errs := make([]error, plSize)
		count := int32(plSize)
		current := int32(0)

		for i, p := range pl {
			pos := i

			p.Then(func(result Result) {
				if result.Err != nil {
					errs[pos] = result.Err
				} else {
					rts[pos] = result.Val
				}

				if atomic.AddInt32(&current, 1) >= count {
					resolve.Return(&Result{
						Val: rts,
						Err: getErr(errs),
					})
				}
			})
		}
	})
}

func getErr(es []error) error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}

	return nil
}
