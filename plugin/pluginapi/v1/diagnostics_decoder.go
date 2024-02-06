package pluginapiv1

import (
	"github.com/hashicorp/hcl/v2"
)

func decodeDiagnosticList(src []*Diagnostic) []*hcl.Diagnostic {
	if src == nil {
		return nil
	}
	dst := make([]*hcl.Diagnostic, len(src))
	for i, v := range src {
		dst[i] = decodeDiagnostic(v)
	}
	return dst
}

func decodeDiagnostic(src *Diagnostic) *hcl.Diagnostic {
	if src == nil {
		return nil
	}
	return &hcl.Diagnostic{
		Severity: decodeDiagnosticSeverity(src.Severity),
		Summary:  src.Summary,
		Detail:   src.Detail,
	}
}

func decodeDiagnosticSeverity(src DiagnosticSeverity) hcl.DiagnosticSeverity {
	switch src {
	case DiagnosticSeverity_DIAGNOSTIC_SEVERITY_ERROR:
		return hcl.DiagError
	case DiagnosticSeverity_DIAGNOSTIC_SEVERITY_WARNING:
		return hcl.DiagWarning
	default:
		return hcl.DiagInvalid
	}
}
