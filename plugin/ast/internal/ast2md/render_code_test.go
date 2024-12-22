package ast2md

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestFindRuns(t *testing.T) {
	assert := assert.New(t)
	for _, tt := range []struct {
		src  string
		char byte
		want []int
	}{
		{"", 'a', []int{}},
		{"a", 'a', []int{1}},
		{"aa", 'a', []int{2}},
		{"aaa", 'a', []int{3}},
		{"bab", 'a', []int{1}},
		{"abab", 'a', []int{1}},
		{"abaabaaab", 'a', []int{1, 2, 3}},
	} {
		t.Run(tt.src, func(t *testing.T) {
			got := findRuns([]byte(tt.src), tt.char)
			gotS := maps.Keys(got)
			slices.Sort(gotS)
			assert.Equal(tt.want, gotS)
		})
	}
}

func TestMakeShortestFence(t *testing.T) {
	assert := assert.New(t)
	for _, tt := range []struct {
		src  string
		char byte
		min  int
		want string
	}{
		{"", '`', 3, "```"},
		{"a", '`', 3, "```"},
		{"`", '`', 3, "```"},
		{"`", '`', 1, "``"},
		{"``", '`', 3, "```"},
		{"``", '`', 1, "`"},
		{"```", '`', 3, "````"},
		{"``` ``", '`', 1, "`"},
		{"``` `` `", '`', 1, "````"},
	} {
		t.Run(tt.src, func(t *testing.T) {
			got := makeShortestFence([]byte(tt.src), tt.char, tt.min)
			assert.Equal(tt.want, string(got))
		})
	}
}
