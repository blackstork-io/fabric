package testtools

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type (
	Assert func(diag *hcl.Diagnostic) bool
)

func severityToString(severity any) string {
	// Yep, that's how it is done
	// https://github.com/hashicorp/hcl/blob/1c5ae8fc88a656ab7bd46da4ff27a20c5a97497b/diagnostic_text.go#L63-L72
	var severityStr string
	switch severity {
	case hcl.DiagError:
		severityStr = "Error"
	case hcl.DiagWarning:
		severityStr = "Warning"
	default:
		// should never happen
		severityStr = "???????"
	}
	return severityStr
}

func IsError(diag *hcl.Diagnostic) bool {
	result := diag.Severity == hcl.DiagError
	if !result {
		slog.Error("Severity assert failed", "expected", severityToString(hcl.DiagError), "value", severityToString(diag.Severity))
	}
	return result
}

func IsWarning(diag *hcl.Diagnostic) bool {
	result := diag.Severity == hcl.DiagWarning
	if !result {
		slog.Error("Severity assert failed", "expected", severityToString(hcl.DiagWarning), "value", severityToString(diag.Severity))
	}
	return result
}

func SummaryContains(substrs ...string) Assert {
	return func(diag *hcl.Diagnostic) bool {
		result := contains(diag.Summary, substrs)
		if !result {
			slog.Error("Summary contains assert failed", "substrings", substrs, "value", diag.Summary)
		}
		return result
	}
}

func DetailContains(substrs ...string) Assert {
	return func(diag *hcl.Diagnostic) bool {
		result := contains(diag.Detail, substrs)
		if !result {
			slog.Error("Detail contains assert failed", "substrings", substrs, "value", diag.Detail)
		}
		return result
	}
}

func SummaryEquals(value string) Assert {
	return func(diag *hcl.Diagnostic) bool {
		result := diag.Summary == value
		if !result {
			slog.Error("Summary equals assert failed", "expected", value, "actual", diag.Summary)
		}
		return result
	}
}

func DetailEquals(value string) Assert {
	return func(diag *hcl.Diagnostic) bool {
		result := diag.Detail == value
		if !result {
			slog.Error("Detail equals assert failed", "expected", value, "actual", diag.Detail)
		}
		return result
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

func CompareDiags[D diagnostics.Generic](t *testing.T, fm map[string]*hcl.File, diags D, asserts [][]Assert) {
	t.Helper()
	compareDiags(t, fm, diagnostics.From(diags), asserts)
}

func compareDiags(t *testing.T, fm map[string]*hcl.File, diags diagnostics.Diag, asserts [][]Assert) {
	t.Helper()
	if !matchBiject(t, diags, asserts) {
		var b strings.Builder
		b.WriteString("Actual diagnostics:\n")
		for _, diag := range diags {
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf("[Severity]: %s\n", severityToString(diag.Severity)))
			b.WriteString(fmt.Sprintf("[Summary]: %s\n", diag.Summary))
			b.WriteString(fmt.Sprintf("[Details]: %s\n", diag.Detail))
			b.WriteString("\n")
		}
		t.Fatal(b.String())
	}
}

func matchBiject(t *testing.T, diags diagnostics.Diag, asserts [][]Assert) bool {
	dgs := []*hcl.Diagnostic(diags)
	if len(dgs) != len(asserts) {
		return false
	}

nextDiag:
	for _, diag := range dgs {
	nextAssertSet:
		for assertSetIdx, assertSet := range asserts {
			for _, a := range assertSet {
				if !a(diag) {
					t.Logf("Assert didn't match the diagnostic")
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
