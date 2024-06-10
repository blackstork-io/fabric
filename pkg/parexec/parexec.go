package parexec

import (
	"log/slog"
	"sync"
	"sync/atomic"
)

var closedChan = func() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}()

// Controls the scheduler. Values: CmdProceed (zero-value) and CmdStop.
type Command uint8

const (
	// Continue async execution (zero-value).
	CmdProceed = Command(iota)
	// Try to stop as soon as possible (start canceling new tasks and stop scheduling already submitted).
	CmdStop
)

// Parallel executor combined with a [sync.Locker] for results.
type Executor[T any] struct {
	idx     atomic.Int64
	tasks   atomic.Int64
	stopC   atomic.Pointer[chan struct{}]
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
	ch := make(chan struct{})
	pe.stopC.Store(&ch)
	pe.cond.L = &sync.Mutex{}
	return pe
}

// Acquires exclusive access to the state captured by the `processor` func of the executor.
// Does _not_ prevent tasks from being scheduled or run, just prevents calls to the processor function.
// Done goroutines would pile up and wait on the lock, so you shouldn't hold this value for long.
func (pe *Executor[T]) Lock() {
	pe.cond.L.Lock()
}

// Releases the exclusive access to the state captured by the `processor` func of the executor.
// Tasks are able to return their results to the processor function again.
func (pe *Executor[T]) Unlock() {
	pe.cond.L.Unlock()
}

// Acquires exclusive access to the state captured by the `processor` func of the executor and
// prevents new tasks from running. Running and finished goroutines would pile up and wait on the lock
// to be able to send their results to the `processor` func of the executor.
// func (pe *Executor[T]) LockAndPause() {
// 	// set limiter to 0
// 	pe.cond.L.Lock()
// }

// Releases the exclusive access to the state captured by the `processor` func of the executor and
// allows new tasks to .
// func (pe *Executor[T]) UnlockAndResume() {
// 	// restore limiter to initial value
// 	pe.cond.L.Unlock()
// }

// Releases the exclusive access to the state captured by the `processor` func of the executor.
// Tasks are able to return their results to the processor function again.
// Additionally resets the stopped status of the executor, making it possible to shedule and run new tasks.
func (pe *Executor[T]) UnlockResume() {
	chPtr := pe.stopC.Load()
	if chPtr == &closedChan {
		ch := make(chan struct{})
		pe.stopC.Store(&ch)
	}
	pe.cond.L.Unlock()
}

// Wait until all scheduled parallel tasks are done (or cancelled), then acquire the lock
// that guarantees exclusive access to the state captured by the `processor` func of the executor.
// Returns whether the executor was stopped (result function returned parexec.CmdStop).
func (pe *Executor[T]) WaitDoneAndLock() (wasStopped bool) {
	pe.cond.L.Lock()
	for pe.tasks.Load() != 0 {
		pe.cond.Wait()
	}
	return pe.isStopped()
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

// must be called under mutex.
func (pe *Executor[T]) stopExecutor() {
	ch := pe.stopC.Swap(&closedChan)
	if ch != &closedChan { // we're the first task that called stopExecutor
		close(*ch)
	}
}

func (pe *Executor[T]) isStopped() bool {
	return pe.stopC.Load() == &closedChan
}

func (pe *Executor[T]) getStopC() <-chan struct{} {
	return *pe.stopC.Load()
}

func (pe *Executor[T]) execute(idx int, fn func() T) {
	limiterActive := pe.limiter != nil
	defer func() {
		pe.taskDone()
		if limiterActive {
			pe.limiter.Return()
		}
		if p := recover(); p != nil {
			slog.Error("Panic in the parallel executor!", "panic:", p)
		}
	}()
	if limiterActive {
		pe.limiter.Take()
	}
	if pe.isStopped() {
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
		pe.stopExecutor()
	}
}

// Execute function f in parallel executor, returns the result into executor's "processor" function.
func (pe *Executor[T]) Go(f func() T) {
	if pe.isStopped() {
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

// Transforms fn (func from (I1, I2) to T) into func(I1, I2) that executes in [Executor] and returns the result
// T into executor's "processor" function.
func GoWithArgs[I1, I2, T any](pe *Executor[T], fn func(I1, I2) T) func(I1, I2) {
	return func(input1 I1, input2 I2) {
		pe.Go(func() T {
			return fn(input1, input2)
		})
	}
}

// Applies func fn to the references of the elements of slice in [Executor] and returns the results
// T into executor's "processor" function.
// Note: when the results of the execution are collected by the processor func of the executor
// indexes of slice items would be contiguous and ordered
// Safety: the array must not be modified until the executor is done!
func MapRef[I, T any](pe *Executor[T], slice []I, fn func(*I) T) {
	if len(slice) == 0 || pe.isStopped() {
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
// indexes of slice items would be contiguous and ordered.
func Map[I, T any](pe *Executor[T], input []I, fn func(I) T) {
	if len(input) == 0 || pe.isStopped() {
		return
	}
	idxStart := pe.taskAdd(len(input))
	for i, input := range input {
		go func(idx int, input I) {
			pe.execute(idx, func() T { return fn(input) })
		}(idxStart+i, input)
	}
}

// Applies func fn to the channel values in [Executor] and returns the results
// T into executor's "processor" function.
// Function returns when the `input` channel is closed or [Executor] is stopped
// Note: when the results of the execution are collected by the processor func of the executor
// indexes of slice items would be ordered, but __not__ contiguous.
func MapChan[I, T any](pe *Executor[T], inputC <-chan I, fn func(I) T) (consumedFully bool) {
	stopC := pe.getStopC()
	for {
		select {
		case <-stopC:
			// executor stopped before inputC was consumed
			return false
		// stopC is not closed, now we can subscribe to inputC
		// otherwise select may choose inputC case even if stopC case is also valid
		default:
			select {
			case <-stopC:
				// executor stopped before inputC was consumed
				return false
			case input, ok := <-inputC:
				if !ok {
					// consumed inputC fully
					return true
				}
				idx := pe.taskAdd(1)
				go func(idx int, input I) {
					pe.execute(idx, func() T { return fn(input) })
				}(idx, input)
			}
		}
	}
}
