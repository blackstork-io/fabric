package backtrace

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type ExecTracer interface {
	FrameEnter(bt *Backtracer) (diag *hcl.Diagnostic)
	FrameExit(bt *Backtracer)
}

// Generates or extends a backtrace
type Backtracer interface {
	AppendBacktrace(*hcl.Diagnostic)
	NewDiagnostic() *hcl.Diagnostic
}

var _ = []Backtracer{
	MessageBacktracer(""),
	(*RangeBacktracer)(nil),
}

// Places a string insted of backtrace location (invisible if an empty string)
type MessageBacktracer string

func (mb MessageBacktracer) NewDiagnostic() *hcl.Diagnostic {
	detail := "Looped back to an object through reference chain:"
	if mb != "" {
		detail = fmt.Sprintf("%s\n  %s", detail, (string)(mb))
	}

	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Circular reference detected",
		Detail:   detail,
		Extra:    BacktraceMarker{},
	}
}

func (mb MessageBacktracer) AppendBacktrace(backtrace *hcl.Diagnostic) {
	if mb == "" {
		return
	}
	backtrace.Detail = fmt.Sprintf(
		"%s\n  %s",
		backtrace.Detail, (string)(mb),
	)
}

type RangeBacktracer hcl.Range

func (rb *RangeBacktracer) AppendBacktrace(backtrace *hcl.Diagnostic) {
	if rb != nil {
		backtrace.Detail = fmt.Sprintf(
			"%s\n  at %s:%d:%d",
			backtrace.Detail, rb.Filename, rb.Start.Line, rb.Start.Column,
		)
	} else {
		backtrace.Detail += "\n  at <missing location info>"
	}
	return
}

func (rb *RangeBacktracer) NewDiagnostic() *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Circular reference detected",
		Detail:   "Looped back to this object through reference chain:",
		Subject:  (*hcl.Range)(rb),
		Extra:    BacktraceMarker{},
	}
}

// Marks the diagnostic as a backtrace
type BacktraceMarker struct{}
