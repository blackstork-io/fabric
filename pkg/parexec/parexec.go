package parexec

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// Controls the scheduler. Values: CmdProceed (zero-value) and CmdStop.
type Command uint8

const (
	// Continue async execution (zero-value).
	CmdProceed = Command(iota)
	// Try to stop as soon as possible (start canceling new tasks and stop scheduling alreday submitted).
	CmdStop
)

// NoLimit doesn't limit the degree of parallel execution.
var NoLimit = (*Limiter)(nil)

// DiskIOLimiter is a shared limiter for all Disk IO tasks.
var DiskIOLimiter = NewLimiter(4) //nolint: gomnd

// CPULimiter is a shared limiter for all CPU-bound tasks.
var CPULimiter = NewLimiter(2 * runtime.NumCPU())

type Limiter struct {
	cond      sync.Cond
	available int
	total     int
}

// Limits parallel executions to at most limit simultaneously.
//
// Can be shared between multiple [ParExecutor]s.
func NewLimiter(limit int) *Limiter {
	l := &Limiter{
		available: limit,
		total:     limit,
	}
	l.cond = *sync.NewCond(&sync.Mutex{})
	return l
}

// Takes a limiter token. Must [Return] it after.
func (l *Limiter) Take() {
	l.cond.L.Lock()
	for l.available <= 0 {
		l.cond.Wait()
	}
	l.available--
	l.cond.L.Unlock()
}

// Returns a token taken with Take.
func (l *Limiter) Return() {
	l.cond.L.Lock()
	if l.available == 0 {
		l.cond.Signal()
	}
	l.available++
	l.cond.L.Unlock()
}

// Parallel executor combined with a [sync.Locker] for results.
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
// 'processor' func is called synchronously (under lock) with result of execution and idx –
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

// Wait unill all scheduled parallel tasks are done (or cancelled), then acquire the lock
// that guarantees exclusive access to the state captured by the `processor` func of the executor.
func (pe *Executor[T]) WaitDoneAndLock() {
	pe.cond.L.Lock()
	for pe.tasks.Load() != 0 {
		pe.cond.Wait()
	}
}

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

func (pe *Executor[T]) execute(idx int, fn func() T) {
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
	res := fn()
	if limiterActive {
		limiterActive = false
		pe.limiter.Return()
	}
	pe.Lock()
	defer pe.Unlock()
	cmd := pe.processor(res, idx)
	if cmd == CmdStop {
		pe.stop.Store(true)
	}
}

// Execute function f in parallel executor, returns the result into executor's "processor" function.
func (pe *Executor[T]) Go(f func() T) {
	if pe.stop.Load() {
		return
	}
	idx := pe.taskAdd(1)
	go pe.execute(idx, f)
}

// Transforms fn (func from I to T) into func(I) that executes in [Executor] and returns the result
// T into executor's "processor" function.
func GoWithArg[I, T any](pe *Executor[T], fn func(I) T) func(I) {
	return func(input I) {
		pe.Go(func() T {
			return fn(input)
		})
	}
}

// Applies func fn to the references of the elements of slice in [Executor] and returns the results
// T into executor's "processor" function.
// Note: when the results of the execution are collected by the processor func of the executor
// indexes of slice items would be contigous and orders
// Safety: the array must not be modified until the executor is done!
func MapRef[I, T any](pe *Executor[T], slice []I, fn func(*I) T) {
	if len(slice) == 0 || pe.stop.Load() {
		return
	}
	idxStart := pe.taskAdd(len(slice))
	for i := range slice {
		go func(idx int, input *I) {
			pe.execute(idx, func() T { return fn(input) })
		}(idxStart+i, &slice[i])
	}
}

// Applies func fn to the copies of the elements of slice in [Executor] and returns the results
// T into executor's "processor" function.
// Note: when the results of the execution are collected by the processor func of the executor
// indexes of slice items would be contigous and orders.
func Map[I, T any](pe *Executor[T], input []I, fn func(I) T) {
	if len(input) == 0 || pe.stop.Load() {
		return
	}
	idxStart := pe.taskAdd(len(input))
	for i, input := range input {
		go func(idx int, input I) {
			pe.execute(idx, func() T { return fn(input) })
		}(idxStart+i, input)
	}
}

// Sets slice[idx] = val, growing slice if needed, and returns the updated slice.
func SetAt[T any](slice []T, idx int, val T) []T {
	needToAlloc := idx - len(slice)
	switch {
	case needToAlloc > 0:
		slice = append(slice, make([]T, needToAlloc)...)

		fallthrough
	case needToAlloc == 0:
		slice = append(slice, val)
	default:
		slice[idx] = val
	}
	return slice
}
