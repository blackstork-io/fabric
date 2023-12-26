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
	pe := parexec.New(parexec.NoLimit, func(res int, idx int) (cmd parexec.Command) {
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
	pe := parexec.New(parexec.NewLimiter(4), func(res int, idx int) (cmd parexec.Command) {
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
