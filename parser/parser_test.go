package parser_test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/parser"
)

func collect[T any](ch <-chan T) (result []T) {
	for v := range ch {
		result = append(result, v)
	}
	return
}

func TestFindFiles(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	fs := fstest.MapFS{
		"f1.fabric":            &fstest.MapFile{},
		"f2.fAbRiC":            &fstest.MapFile{},
		"f3.not_fabric":        &fstest.MapFile{},
		"subdir/f4.fAbRiC":     &fstest.MapFile{},
		"subdir/f5.not_fabric": &fstest.MapFile{},
	}

	type testCase struct {
		desc      string
		recursive bool
		expected  []string
	}

	testCases := []testCase{
		{
			desc:      "Recursive",
			recursive: true,
			expected: []string{
				"f1.fabric",
				"f2.fAbRiC",
				"subdir/f4.fAbRiC",
			},
		},
		{
			desc:      "Non-recursive",
			recursive: false,
			expected: []string{
				"f1.fabric",
				"f2.fAbRiC",
			},
		},
	}
	for _, tC := range testCases {
		func(tC testCase) {
			t.Run(tC.desc, func(t *testing.T) {
				t.Parallel()
				var res []string

				diags := parser.FindFabricFiles(fs, tC.recursive, func(path string) {
					res = append(res, path)
				})

				assert.Equal(
					tC.expected,
					res,
				)
				assert.Empty(diags)
			})
		}(tC)
	}
}
