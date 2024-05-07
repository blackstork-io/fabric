package diagtest

import (
	"bytes"
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

type Asserts [][]Assert

func (asserts Asserts) AssertMatch(tb testing.TB, diags []*hcl.Diagnostic, fileMap map[string]*hcl.File) {
	tb.Helper()
	if matchBiject(diags, asserts) {
		return
	}

	var buf strings.Builder
	buf.WriteString("Expected ")
	switch len(asserts) {
	case 0:
		buf.WriteString("no diagnostics\n")
	case 1:
		buf.WriteString("1 diagnostic:\n")
	default:
		fmt.Fprintf(&buf, "%d diagnostics:\n", len(asserts))
	}
	for _, assertSet := range asserts {
		fmt.Fprintln(&buf, "{")
		for _, assert := range assertSet {
			fmt.Fprintf(&buf, "    %s\n", assert)
		}
		fmt.Fprintln(&buf, "},")
	}
	buf.WriteString("\nGot ")
	switch len(diags) {
	case 0:
		buf.WriteString("no diagnostics\n")
	case 1:
		buf.WriteString("1 diagnostic:\n")
	default:
		fmt.Fprintf(&buf, "%d diagnostics:\n", len(diags))
	}
	if len(diags) > 0 {
		diagnostics.PrintDiags(&buf, diags, fileMap, false)
	}
	tb.Fatalf("\n\n%s", buf.String())
}

func sliceRemove[T any](s []T, pos int) []T {
	var tmp T
	s[pos] = s[len(s)-1]
	s[len(s)-1] = tmp
	return s[:len(s)-1]
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

func AssertNoErrors(tb testing.TB, diags []*hcl.Diagnostic, fileMap map[string]*hcl.File, msgs ...any) {
	tb.Helper()
	if len(diags) == 0 {
		return
	}
	var buf bytes.Buffer
	diagnostics.PrintDiags(&buf, diags, fileMap, false)
	msgs = append(msgs, buf.String())
	if hcl.Diagnostics(diags).HasErrors() {
		tb.Error(msgs...)
	} else {
		tb.Log(msgs...)
	}
}
