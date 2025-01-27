package nodes

import (
	"cmp"
)

type Path []int

// Clone returns a copy of the path.
func (p Path) Clone() Path {
	clone := make(Path, len(p))
	copy(clone, p)
	return clone
}

func (p Path) CommonPrefixLen(other Path) int {
	minLen := min(len(p), len(other))
	// bounds check elimination
	if minLen > len(p) || minLen > len(other) {
		panic("unreachable")
	}
	for i := 0; i < minLen; i++ {
		if p[i] != other[i] {
			return i
		}
	}
	return minLen
}

// ComparePaths compares two paths. It returns -1 if a < b, 0 if a == b, and 1 if a > b.
// Comparison is performed element-wise. If all elements are equal, the shorter path
// is considered greater (as it refers to the higher-level node in AST).
func ComparePaths(a, b Path) int {
	minLen := min(len(a), len(b))
	// bounds check elimination
	if minLen > len(a) || minLen > len(b) {
		panic("unreachable")
	}
	for i := 0; i < minLen; i++ {
		c := cmp.Compare(a[i], b[i])
		if c != 0 {
			return c
		}
	}
	return cmp.Compare(len(b), len(a))
}
