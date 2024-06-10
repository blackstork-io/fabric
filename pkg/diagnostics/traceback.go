package diagnostics

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type TracebackExtra = *tracebackExtra

type tracebackExtra struct {
	Traceback []*hcl.Range
}

func (tb *tracebackExtra) improveDiagnostic(diag *hcl.Diagnostic) {
	sb := []byte(diag.Detail)
	for _, rng := range tb.Traceback {
		if rng != nil {
			sb = fmt.Appendf(sb, "\n  at %s:%d:%d",
				rng.Filename, rng.Start.Line, rng.Start.Column,
			)
		} else {
			sb = append(sb, "\n  at <missing location info>"...)
		}
	}
	diag.Detail = string(sb)
}

func NewTracebackExtra() TracebackExtra {
	return &tracebackExtra{}
}
