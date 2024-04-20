package testtools

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
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

func dumpDiag(diag *hcl.Diagnostic) string {
	var sev string
	switch diag.Severity {
	case hcl.DiagError:
		sev = "hcl.DiagError"
	case hcl.DiagWarning:
		sev = "hcl.DiagWarning"
	default:
		sev = fmt.Sprintf("hcl.DiagInvalid(%d)", diag.Severity)
	}
	return fmt.Sprintf("Severity: %s; Summary: %q; Detail: %q", sev, diag.Summary, diag.Detail)
}

func dumpDiags(diags diagnostics.Diag) string {
	if len(diags) == 0 {
		return "no diagnostics"
	}
	return strings.Join(utils.FnMap(diags, dumpDiag), "\n")
}

func CompareDiags[D diagnostics.Generic](t *testing.T, fm map[string]*hcl.File, diags D, asserts [][]Assert) {
	t.Helper()
	compareDiags(t, fm, diagnostics.From(diags), asserts)
}

func compareDiags(t *testing.T, fm map[string]*hcl.File, diags diagnostics.Diag, asserts [][]Assert) {
	t.Helper()
	if !matchBiject(diags, asserts) {
		var buf strings.Builder
		diagnostics.PrintDiags(&buf, diags, fm, false)
		t.Fatalf("\n\n%s", buf.String())
	}
}

func matchBiject(diags diagnostics.Diag, asserts [][]Assert) bool {
	dgs := []*hcl.Diagnostic(diags)
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
