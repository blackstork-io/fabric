package once

import (
	"sync"
)

type onceValues[T1, T2 any] struct {
	status     status
	res1       T1
	res2       T2
	fn         func() (T1, T2)
	panicValue any
	mu         sync.Mutex
}

func (o *onceValues[T1, T2]) do() (T1, T2) {
	status := o.status.AtomicLoad()
	if !status.IsSuccess() {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(status)
	}
	return o.res1, o.res2
}

func (o *onceValues[T1, T2]) doSlow(status status) {
	// Pretty sure that this could be inlined in the do, but that's too much optimization
	// for the (hopefully) rare "panic" case
	// This check can also occur under mutex, but doing it this way seems just a bit neater
	if status.IsPanic() {
		panic(o.panicValue)
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	switch {
	case o.status.IsUnexecuted():
		defer func() {
			// we performed our call, allow GC to collect the fn (and everything captured by the fn)
			o.fn = nil
			if !status.IsSuccess() {
				// Haven't set to success, this means we've panicked
				o.panicValue = recover()
				o.status.AtomicStore(statusPanic)
				panic(o.panicValue)
			}
		}()
		o.res1, o.res2 = o.fn()
		// using status input variable to avoid touching o.status in defer above atomically (more expensive)
		// or non-atomically (the guarantees are not *that* clear)
		status = statusSuccess
		o.status.AtomicStore(statusSuccess)
	case o.status.IsPanic():
		panic(o.panicValue)
	case o.status.IsSuccess():
		return
	}
}

func Values[T1, T2 any](f func() (T1, T2)) func() (T1, T2) {
	return (&onceValues[T1, T2]{
		fn:     f,
		status: statusUnexecuted,
	}).do
}
