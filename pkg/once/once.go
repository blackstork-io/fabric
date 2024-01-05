package once

// Temporarily using sync's versions until we add a full test suite for the custom "ones"
// or decide that we don't want to allow GC of the captured args

import "sync"

func Func(f func()) func() {
	return sync.OnceFunc(f)
}

func Value[T any](f func() T) func() T {
	return sync.OnceValue(f)
}

func Values[T1, T2 any](f func() (T1, T2)) func() (T1, T2) {
	return sync.OnceValues(f)
}
