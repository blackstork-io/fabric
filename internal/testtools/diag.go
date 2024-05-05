package testtools

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
)

type Assert interface {
	Assert(diag *hcl.Diagnostic) bool
	fmt.Stringer
}

type severityAssert struct {
	severity hcl.DiagnosticSeverity
}

func (s severityAssert) Assert(diag *hcl.Diagnostic) bool {
	return diag.Severity == s.severity
}

// String implements Assert.
func (s severityAssert) String() string {
	switch s.severity {
	case hcl.DiagError:
		return "An error"
	case hcl.DiagWarning:
		return "A warning"
	default:
		return fmt.Sprintf("Severity to be equal to hcl.DiagnosticSeverity(%d)", s.severity)
	}
}

var _ Assert = (*severityAssert)(nil)

var (
	IsError   = severityAssert{hcl.DiagError}
	IsWarning = severityAssert{hcl.DiagWarning}
)

type containsAssert struct {
	isSummary bool // if false - detail
	substrs   []string
}

func (c *containsAssert) Assert(diag *hcl.Diagnostic) bool {
	var str string
	if c.isSummary {
		str = diag.Summary
	} else {
		str = diag.Detail
	}
	str = strings.ToLower(str)
	for _, substr := range c.substrs {
		if !strings.Contains(str, strings.ToLower(substr)) {
			return false
		}
	}
	return true
}

// String implements Assert.
func (c *containsAssert) String() string {
	var attrName string
	if c.isSummary {
		attrName = "Summary"
	} else {
		attrName = "Detail"
	}

	return fmt.Sprintf(
		"%s to contain: %s",
		attrName,
		strings.Join(
			utils.FnMap(
				c.substrs,
				func(s string) string {
					return fmt.Sprintf("%q", s)
				},
			),
			", ",
		),
	)
}

func SummaryContains(substrs ...string) Assert {
	return &containsAssert{
		isSummary: true,
		substrs:   substrs,
	}
}

func DetailContains(substrs ...string) Assert {
	return &containsAssert{
		isSummary: false,
		substrs:   substrs,
	}
}

func sliceRemove[T any](s []T, pos int) []T {
	var tmp T
	s[pos] = s[len(s)-1]
	s[len(s)-1] = tmp
	return s[:len(s)-1]
}

func CompareDiags[D diagnostics.Generic](t *testing.T, fm map[string]*hcl.File, diags D, asserts [][]Assert) {
	t.Helper()
	compareDiags(t, fm, diagnostics.From(diags), asserts)
}

func compareDiags(t *testing.T, fm map[string]*hcl.File, diags diagnostics.Diag, asserts [][]Assert) {
	t.Helper()
	if !matchBiject(diags, asserts) {
		var buf strings.Builder
		buf.WriteString("Expected ")
		if len(asserts) == 0 {
			buf.WriteString("no diagnostics\n")
		} else if len(asserts) == 1 {
			buf.WriteString("1 diagnostic:\n")
		} else {
			fmt.Fprintf(&buf, "%d diagnostics:\n", len(asserts))
		}
		for _, assertSet := range asserts {
			fmt.Fprintln(&buf, "{")
			for _, assert := range assertSet {
				fmt.Fprintf(&buf, "    %s\n", assert)
			}
			fmt.Fprintln(&buf, "},")
		}
		buf.WriteString("\nGot:\n\n")
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
			if len(assertSet) == 0 {
				panic("assert set has length 0")
			}
			for _, assert := range assertSet {
				if !assert.Assert(diag) {
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
