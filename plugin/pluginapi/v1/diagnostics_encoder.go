package pluginapiv1

import (
	"github.com/hashicorp/hcl/v2"
)

func encodeDiagnosticList(src []*hcl.Diagnostic) []*Diagnostic {
	if src == nil {
		return nil
	}
	dst := make([]*Diagnostic, len(src))
	for i, v := range src {
		dst[i] = encodeDiagnostic(v)
	}
	return dst
}

func encodeDiagnostic(src *hcl.Diagnostic) *Diagnostic {
	if src == nil {
		return nil
	}
	return &Diagnostic{
		Severity: encodeDiagnosticSeverity(src.Severity),
		Summary:  src.Summary,
		Detail:   src.Detail,
	}
}

func encodeDiagnosticSeverity(src hcl.DiagnosticSeverity) DiagnosticSeverity {
	switch src {
	case hcl.DiagError:
		return DiagnosticSeverity_DIAGNOSTIC_SEVERITY_ERROR
	case hcl.DiagWarning:
		return DiagnosticSeverity_DIAGNOSTIC_SEVERITY_WARNING
	default:
		return DiagnosticSeverity_DIAGNOSTIC_SEVERITY_UNSPECIFIED
	}
}
