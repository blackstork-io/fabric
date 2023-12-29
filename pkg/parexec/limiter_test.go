package parexec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimiterOnPanicingJob(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Log("panic caught!", r)
		}
	}()

	t.Parallel()
	assert := assert.New(t)

	limiter := NewLimiter(4)

	var resArr []int

	pe := New(limiter, func(res int, idx int) (cmd Command) {
		resArr = append(resArr, res)
		return
	})

	for i := 0; i < 8; i++ {
		pe.Go(func() int {
			panic("panicErr")
			// return 0
		})
	}

	wasStopped := pe.WaitDoneAndLock()
	assert.False(wasStopped)

	limiter.cond.L.Lock()
	assert.Equal(limiter.total, limiter.available)
	limiter.cond.L.Unlock()

	pe.UnlockResume()
	// check that everything is still usable after panic in the executor
	for i := 0; i < 8; i++ {
		pe.Go(func() int {
			return 1
		})
	}

	wasStopped = pe.WaitDoneAndLock()

	assert.False(wasStopped)
	limiter.cond.L.Lock()
	assert.Equal(limiter.total, limiter.available)
	limiter.cond.L.Unlock()

	assert.Len(resArr, 8)

	for _, val := range resArr {
		assert.Equal(1, val)
	}
}
