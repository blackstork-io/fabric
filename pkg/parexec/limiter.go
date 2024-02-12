package parexec

import (
	"runtime"
	"sync"
)

// NoLimit doesn't limit the degree of parallel execution.
var NoLimit = (*Limiter)(nil)

// DiskIOLimiter is a shared limiter for all Disk IO tasks.
var DiskIOLimiter = NewLimiter(4)

// CPULimiter is a shared limiter for all CPU-bound tasks.
var CPULimiter = NewLimiter(2 * runtime.NumCPU())

// 4 -> 4+2
// 8 -> 8+3
// 16 -> 16+4

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
	l.available++
	l.cond.L.Unlock()
	l.cond.Signal()
}
