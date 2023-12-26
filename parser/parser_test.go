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
		expected  []parser.FindFilesResult
	}

	testCases := []testCase{
		{
			desc:      "Recursive",
			recursive: true,
			expected: []parser.FindFilesResult{
				{Path: "f1.fabric"},
				{Path: "f2.fAbRiC"},
				{Path: "subdir/f4.fAbRiC"},
			},
		},
		{
			desc:      "Non-recursive",
			recursive: false,
			expected: []parser.FindFilesResult{
				{Path: "f1.fabric"},
				{Path: "f2.fAbRiC"},
			},
		},
	}
	for _, tC := range testCases {
		func(tC testCase) {
			t.Run(tC.desc, func(t *testing.T) {
				t.Parallel()
				assert.Equal(
					tC.expected,
					collect(parser.FindFiles(fs, tC.recursive)),
				)
			})
		}(tC)
	}
}
