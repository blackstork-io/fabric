package diag_test

import (
	"strings"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type (
	Assert func(diag *hcl.Diagnostic) bool
)

func IsError(diag *hcl.Diagnostic) bool {
	return diag.Severity == hcl.DiagError
}

func IsWarning(diag *hcl.Diagnostic) bool {
	return diag.Severity == hcl.DiagWarning
}

func SummaryContains(substrs ...string) Assert {
	return func(diag *hcl.Diagnostic) bool {
		return contains(diag.Summary, substrs)
	}
}

func DetailContains(substrs ...string) Assert {
	return func(diag *hcl.Diagnostic) bool {
		return contains(diag.Detail, substrs)
	}
}

func contains(str string, substrs []string) bool {
	str = strings.ToLower(str)
	for _, substr := range substrs {
		if !strings.Contains(str, strings.ToLower(substr)) {
			return false
		}
	}
	return true
}

func sliceRemove[T any](s []T, pos int) []T {
	var tmp T
	s[pos] = s[len(s)-1]
	s[len(s)-1] = tmp
	return s[:len(s)-1]
}

func MatchBiject(diags diagnostics.Diag, asserts [][]Assert) bool {
	dgs := ([]*hcl.Diagnostic)(diags)
	if len(dgs) != len(asserts) {
		return false
	}

nextDiag:
	for _, diag := range dgs {
	nextAssertSet:
		for assertSetIdx, assertSet := range asserts {
			for _, assert := range assertSet {
				if !assert(diag) {
					continue nextAssertSet
				}
			}
			// all asserts in assert set have matched, remove
			asserts = sliceRemove(asserts, assertSetIdx)
			continue nextDiag
		}
		// can't find assert set matching this diag
		return false
	}

	return true
}
