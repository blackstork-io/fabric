package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyFn(t *testing.T) {
	t.Parallel()
	product := 1
	ApplyFn([]int{2, 3, 7}, func(i int) {
		product *= i
	})
	assert.Equal(t, 42, product)
}

func TestApplyFnRef(t *testing.T) {
	t.Parallel()
	product := 1
	ApplyFnRef([]int{2, 3, 7}, func(i *int) {
		product *= *i
	})
	assert.Equal(t, 42, product)
}

// func ApplyFn[T any](slice []T, fn func(T)) {
// 	for _, it := range slice {
// 		fn(it)
// 	}
// }

// func ApplyFnRef[T any](slice []T, fn func(*T)) {
// 	for _, it := range slice {
// 		fn(&it)
// 	}
// }
