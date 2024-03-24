package cachedval_test

// import (
// 	"sync"
// 	"testing"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/blackstork-io/fabric/pkg/backtrace"
// 	"github.com/blackstork-io/fabric/pkg/cachedval"
// 	"github.com/blackstork-io/fabric/pkg/diagnostics"
// )

// type barrier struct {
// 	wg sync.WaitGroup
// }

// func (b *barrier) Wait() {
// 	b.wg.Done()
// 	b.wg.Wait()
// }

// func NewBarrier(n int) *barrier {
// 	var b barrier
// 	b.wg.Add(n)
// 	return &b
// }

// func TestCachedVal(t *testing.T) {
// 	t.Parallel()
// 	assert := assert.New(t)

// 	crd := cachedval.NewCRD()

// 	bt0 := backtrace.MessageBacktracer("bt0")
// 	bt1 := backtrace.MessageBacktracer("bt1")
// 	bt2 := backtrace.MessageBacktracer("bt2")

// 	v0 := cachedval.New[int](bt0)
// 	v1 := cachedval.New[int](bt1)
// 	v2 := cachedval.New[int](bt2)

// 	_, diags := v0.Get(crd, func() (int, diagnostics.Diag) {
// 		return v1.Get(crd, func() (int, diagnostics.Diag) {
// 			return v2.Get(crd, func() (int, diagnostics.Diag) {
// 				return v1.Get(crd, func() (int, diagnostics.Diag) {
// 					assert.Fail("Should not enter ref loop")
// 					return 0, nil
// 				})
// 			})
// 		})
// 	})
// 	assert.NotEmpty(diags)
// 	assert.Len(diags, 1)
// 	assert.Equal(diags[0].Detail, `Looped back to an object through reference chain:
//   bt1
//   bt2
//   bt1
// Reference loop entered because of:
//   bt0`)
// }

// func TestCachedValAccess(t *testing.T) {
// 	t.Parallel()
// 	assert := assert.New(t)

// 	bt0 := backtrace.MessageBacktracer("bt0")

// 	v0 := cachedval.New[*int](bt0)
// 	call_count := 0

// 	vals := make(chan *int)
// 	var wg sync.WaitGroup

// 	const NumRunners = 5

// 	wg.Add(NumRunners)
// 	for i := 0; i < NumRunners; i++ {
// 		go func(call_count *int) {
// 			defer wg.Done()
// 			crd := cachedval.NewCRD()
// 			n, diags := v0.Get(crd, func() (val *int, diags diagnostics.Diag) {
// 				*call_count += 1
// 				v := 123
// 				val = &v
// 				return
// 			})
// 			assert.Empty(diags)
// 			vals <- n
// 		}(&call_count)
// 	}
// 	go func() {
// 		wg.Wait()
// 		close(vals)
// 	}()

// 	first_val := <-vals
// 	for val := range vals {
// 		assert.Same(first_val, val)
// 	}
// 	assert.Equal(1, call_count)
// }
