package plugin

type DiagnosticSeverity int

const (
	DiagInvalid DiagnosticSeverity = iota
	DiagError
	DiagWarning
)

type Diagnostic struct {
	Severity DiagnosticSeverity
	Summary  string
	Detail   string
}

type Diagnostics []*Diagnostic
