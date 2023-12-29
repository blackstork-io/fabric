package parexec_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/pkg/parexec"
)

func TestSetAt(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		idx   int
		val   int
	}{
		{
			name:  "Regular append",
			slice: []int{0, 1, 2},
			idx:   3,
			val:   1337,
		},
		{
			name:  "Set existing",
			slice: []int{0, 1, 2},
			idx:   1,
			val:   1337,
		},
		{
			name:  "Set and extend",
			slice: []int{0, 1, 2},
			idx:   10,
			val:   1337,
		},
	}
	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			assert := assert.New(t)
			orig := slices.Clone(tC.slice)
			res := parexec.SetAt(tC.slice, tC.idx, tC.val)

			assert.True(len(res) == max(len(orig), tC.idx+1))

			for i := range res {
				switch {
				case i == tC.idx: // value at `idx` should become `val`
					assert.Equal(tC.val, res[i])
				case i < len(orig): // other values from the original slice shouldn't change
					assert.Equal(orig[i], res[i])
				default: // i >= len(orig)
					// other values should be filled by the zero value
					assert.Zero(res[i])
				}
			}
		})
	}
}
