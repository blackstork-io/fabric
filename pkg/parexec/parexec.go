package parexec

import (
	"sync"
	"sync/atomic"
)

// Controls the scheduler. Values: PROCEED (zero-value) and STOP
type Command uint8

const (
	// Continue async execution (zero-value)
	PROCEED = Command(iota)
	// Try to stop as soon as possible (start canceling new tasks and stop scheduling alreday submitted)
	STOP
)

// Sets no limit to the degree of parallel execution
var NO_LIMIT = (*Limiter)(nil)

type Limiter struct {
	cond      sync.Cond
	availible int
	total     int
}

// Limits parallel executions to at most limit simultaneously.
//
// Can be shared between multiple [ParExecutor]s
func NewLimiter(limit int) *Limiter {
	l := &Limiter{
		availible: limit,
		total:     limit,
	}
	l.cond = *sync.NewCond(&sync.Mutex{})
	return l
}

// Takes a limiter token. Must [Return] it after
func (l *Limiter) Take() {
	l.cond.L.Lock()
	for l.availible <= 0 {
		l.cond.Wait()
	}
	l.availible--
	l.cond.L.Unlock()
}

// Returns a token taken with Take
func (l *Limiter) Return() {
	l.cond.L.Lock()
	if l.availible == 0 {
		l.cond.Signal()
	}
	l.availible++
	l.cond.L.Unlock()
}

// Parallel executor combined with a [sync.Locker] for results
type Executor[T any] struct {
	idx     atomic.Int64
	tasks   atomic.Int64
	stop    atomic.Bool
	limiter *Limiter

	cond      sync.Cond
	processor func(res T, idx int) Command
}

// Create a new parallel executor
//
// 'processor' func is called syncronously (under lock) with result of execution and idx –
// a monotonically increasing from 0 number, reflecting the order in which the tasks were scheduled
//
// ParExecutor is also a mutex around data captured by the "processor" closure as soon as it's created.
// To be safe, use WaitDoneAndLock() to access this data.
func New[T any](limiter *Limiter, processor func(res T, idx int) Command) *Executor[T] {
	pe := &Executor[T]{
		processor: processor,
		limiter:   limiter,
	}
	pe.cond.L = &sync.Mutex{}
	return pe
}

func (pe *Executor[T]) Lock() {
	pe.cond.L.Lock()
}

func (pe *Executor[T]) Unlock() {
	pe.cond.L.Unlock()
}

// Wait unill all scheduled parallel tasks are done (or cancelled)
// After this function returns and before any other tasks are scheduled nothing
// will execute processor function and you can access its data

// Allows scheduled tasks to modify the "processor" data, wait untill all tasks are done and return
func (pe *Executor[T]) WaitDoneAndLock() {
	pe.cond.L.Lock()
	for pe.tasks.Load() != 0 {
		pe.cond.Wait()
	}
}

// n >= 1!
func (pe *Executor[T]) taskAdd(n int) (idx int) {
	if n < 1 {
		panic("n must be strictly positive")
	}
	n64 := int64(n)
	pe.tasks.Add(n64)
	return int(pe.idx.Add(n64) - n64)
}

func (pe *Executor[T]) taskDone() {
	taskCount := pe.tasks.Add(-1)
	if taskCount == 0 {
		pe.cond.Broadcast()
	}
}

func (pe *Executor[T]) goroutineBody(idx int, f func() T) {
	defer pe.taskDone()
	limiterActive := pe.limiter != nil
	if limiterActive {
		pe.limiter.Take()
		defer func() {
			if limiterActive {
				pe.limiter.Return()
			}
		}()
	}
	if pe.stop.Load() {
		return
	}
	res := f()
	if limiterActive {
		limiterActive = false
		pe.limiter.Return()
	}
	pe.Lock()
	defer pe.Unlock()
	cmd := pe.processor(res, idx)
	if cmd == STOP {
		pe.stop.Store(true)
	}
}

// Execute function f in parallel executor, returns the result into executor's "processor" function
func (pe *Executor[T]) Go(f func() T) {
	if pe.stop.Load() {
		return
	}
	idx := pe.taskAdd(1)
	go pe.goroutineBody(idx, f)
}

func GoWithArg[I, T any](pe *Executor[T], f func(I) T) func(I) {
	return func(input I) {
		pe.Go(func() T {
			return f(input)
		})
	}
}

func MapRef[I, T any](pe *Executor[T], input []I, f func(*I) T) {
	if len(input) == 0 || pe.stop.Load() {
		return
	}
	idxStart := pe.taskAdd(len(input))
	for i := range input {
		go func(idx int, input *I) {
			pe.goroutineBody(idx, func() T { return f(input) })
		}(idxStart+i, &input[i])
	}
}

func Map[I, T any](pe *Executor[T], input []I, f func(I) T) {
	if len(input) == 0 || pe.stop.Load() {
		return
	}
	idxStart := pe.taskAdd(len(input))
	for i, input := range input {
		go func(idx int, input I) {
			pe.goroutineBody(idx, func() T { return f(input) })
		}(idxStart+i, input)
	}
}

// Sets s[idx] = val, growing s if needed, and returns updated slice
func SetAt[T any](s []T, idx int, val T) []T {
	needToAlloc := idx - len(s)
	switch {
	case needToAlloc > 0:
		s = append(s, make([]T, needToAlloc)...)
		fallthrough
	case needToAlloc == 0:
		s = append(s, val)
	default:
		s[idx] = val
	}
	return s
}
