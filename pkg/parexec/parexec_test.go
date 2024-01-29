package parexec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/pkg/parexec"
)

const Len = 2048

func TestBasic(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resArr := make([]int, Len)
	pe := parexec.New(parexec.NoLimit, func(res, idx int) (cmd parexec.Command) {
		resArr[idx] = res
		return
	})

	fn := parexec.GoWithArg(pe, func(i int) int {
		return ^i
	})

	for i := range resArr {
		fn(i)
	}
	pe.WaitDoneAndLock()
	for i := range resArr {
		assert.Exactly(resArr[i], ^i)
	}
}

func TestReenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resArr := make([]int, Len*3)
	pe := parexec.New(parexec.NewLimiter(4), func(res, idx int) (cmd parexec.Command) {
		resArr[idx] = res
		return
	})

	fn := parexec.GoWithArg(pe, func(i int) int {
		return ^i
	})

	for i := 0; i < Len; i++ {
		fn(i)
	}
	pe.WaitDoneAndLock()
	var i int
	for i = 0; i < Len; i++ {
		assert.Exactly(^(i % Len), resArr[i])
	}
	for ; i < len(resArr); i++ {
		assert.Exactly(0, resArr[i])
	}
	for i := 0; i < Len; i++ {
		fn(i)
	}
	// tasks are launched but can't write to the processor: we are holding the lock
	for i = 0; i < Len; i++ {
		assert.Exactly(^(i % Len), resArr[i])
	}
	for ; i < len(resArr); i++ {
		assert.Exactly(0, resArr[i])
	}
	pe.Unlock()
	pe.WaitDoneAndLock()
	for i = 0; i < Len*2; i++ {
		assert.Exactly(^(i % Len), resArr[i])
	}
	for ; i < len(resArr); i++ {
		assert.Exactly(0, resArr[i])
	}
	pe.Unlock()
	for i := 0; i < Len; i++ {
		fn(i)
	}
	pe.WaitDoneAndLock()
	for i = 0; i < Len*3; i++ {
		assert.Exactly(^(i % Len), resArr[i])
	}
}

func TestCancel(t *testing.T) {
	const jobCount = 8
	t.Parallel()
	assert := assert.New(t)

	var resArr []int
	pe := parexec.New(parexec.NoLimit, func(res, idx int) (cmd parexec.Command) {
		resArr = append(resArr, res)
		if len(resArr) == 1 {
			return parexec.CmdStop
		}
		// A single stop command should be enough
		return
	})

	for i := 0; i < jobCount; i++ {
		pe.Go(func() int {
			return 0
		})
	}
	// max 8 submitted
	wasStopped := pe.WaitDoneAndLock()
	assert.True(wasStopped)
	jobsDone := len(resArr)
	assert.True(jobsDone >= 1 && jobsDone <= 8)

	pe.Unlock() // did not clear the "stopped" status, so next submitted jobs are silently dropped

	for i := 0; i < jobCount; i++ {
		pe.Go(func() int {
			return 1
		})
	}

	wasStopped = pe.WaitDoneAndLock()
	assert.True(wasStopped, "must be still stopped")
	assert.Len(resArr, jobsDone, "shouldn't execute anything: executor stopped")

	pe.UnlockResume()

	for i := 0; i < jobCount; i++ {
		pe.Go(func() int {
			return 2
		})
	}

	wasStopped = pe.WaitDoneAndLock()
	assert.False(wasStopped, "resume should've resumed execution")
	assert.Len(resArr, jobsDone+jobCount, "all jobCount jobs should've been executed")

	for i := 0; i < jobsDone; i++ {
		assert.Equal(0, resArr[i])
	}
	for i := jobsDone; i < jobsDone+jobCount; i++ {
		assert.Equal(2, resArr[i])
	}
}
