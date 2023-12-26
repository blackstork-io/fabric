// Tests primarily desingned to reveal race conditions, advised to run with -race

package parexec_test

import (
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/pkg/parexec"
)

func TestManyTasks(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resArr := make([]int, Len*64)
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

func TestLimiter(t *testing.T) {
	t.Parallel()
	for _, limit := range []int{1, 2, 3, 4, 8, 16, 1024} {
		func(limit int) {
			t.Run(strconv.Itoa(limit), func(t *testing.T) {
				t.Parallel()
				assert := assert.New(t)
				resArr := make([]int, Len)
				pe := parexec.New(parexec.NewLimiter(limit), func(res int, idx int) (cmd parexec.Command) {
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
					assert.Exactly(^i, resArr[i])
				}
			})
		}(limit)
	}
}

func TestManyTaskGivers(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	type result struct {
		workerID int
		num      int
	}

	resMap := map[int][]int{}
	pe := parexec.New(parexec.NoLimit, func(res result, idx int) (cmd parexec.Command) {
		resMap[res.workerID] = append(resMap[res.workerID], res.num)
		return
	})

	const workers = 8
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			defer wg.Done()
			f := parexec.GoWithArg(pe, func(i int) result {
				return result{
					workerID: worker,
					num:      i,
				}
			})
			for i := 0; i < Len; i++ {
				f(i)
			}
		}(worker)
	}
	wg.Wait()
	pe.WaitDoneAndLock()
	assert.Len(resMap, workers)
	for worker := 0; worker < workers; worker++ {
		assert.Contains(resMap, worker)
		assert.Len(resMap[worker], Len)
		slices.Sort(resMap[worker])
		for i := 0; i < Len; i++ {
			assert.Equal(i, resMap[worker][i])
		}
	}
}

func TestManyTaskGiversAndWaiters(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	type result struct {
		workerID int
		num      int
	}

	resMap := map[int][]int{}
	pe := parexec.New(parexec.NoLimit, func(res result, idx int) (cmd parexec.Command) {
		resMap[res.workerID] = append(resMap[res.workerID], res.num)
		return
	})

	const workers = 8
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			defer wg.Done()
			f := parexec.GoWithArg(pe, func(i int) result {
				time.Sleep(time.Duration(2*worker) * time.Millisecond)
				return result{
					workerID: worker,
					num:      i,
				}
			})
			for i := 0; i < Len; i++ {
				f(i)
			}
			pe.WaitDoneAndLock()
			// we're guaranteed that all submitted work for this worker is done. Other
			// workers my be waiting right now
			defer pe.Unlock()
			assert.Contains(resMap, worker)
			assert.Len(resMap[worker], Len)
			slices.Sort(resMap[worker])
			for i := 0; i < Len; i++ {
				assert.Equal(i, resMap[worker][i])
			}
		}(worker)
	}
	wg.Wait()
	pe.WaitDoneAndLock()
	assert.Len(resMap, workers)
	for worker := 0; worker < workers; worker++ {
		assert.Contains(resMap, worker)
		assert.Len(resMap[worker], Len)
		slices.Sort(resMap[worker])
		for i := 0; i < Len; i++ {
			assert.Equal(i, resMap[worker][i])
		}
	}
}
