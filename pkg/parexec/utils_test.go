package parexec_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/pkg/parexec"
	"github.com/blackstork-io/fabric/pkg/utils"
)

func TestSetAt(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name  string
		slice []int
		idx   int
		val   int
	}

	tests := []testCase{
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

	utils.ApplyFn(tests, func(tc testCase) {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			orig := slices.Clone(tc.slice)
			res := parexec.SetAt(tc.slice, tc.idx, tc.val)

			assert.Len(res, max(len(orig), tc.idx+1))

			for i := range res {
				switch {
				case i == tc.idx: // value at `idx` should become `val`
					assert.Equal(tc.val, res[i])
				case i < len(orig): // other values from the original slice shouldn't change
					assert.Equal(orig[i], res[i])
				default: // i >= len(orig)
					// other values should be filled by the zero value
					assert.Zero(res[i])
				}
			}
		})
	})
}
