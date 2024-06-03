package diagnostics

// More ergonomic wrapper over hcl.Diagnostics.

import (
	"github.com/hashicorp/hcl/v2"
)

type extraList []any

func (e extraList) extraListSigil() {}

func diagnosticExtra[T any](extra any) (_ T, _ bool) {
	for extra != nil {
		switch extraT := extra.(type) {
		case T:
			return extraT, true
		case extraList:
			for _, extra := range extraT {
				val, found := diagnosticExtra[T](extra)
				if found {
					return val, found
				}
			}
		case hcl.DiagnosticExtraUnwrapper:
			extra = extraT.UnwrapDiagnosticExtra()
			continue
		}
		break
	}
	return
}

func DiagnosticExtra[T any](diag *hcl.Diagnostic) (_ T, _ bool) {
	return diagnosticExtra[T](diag.Extra)
}

func DiagnosticsExtra[T any](diags Diag) (val T, found bool) {
	for _, diag := range diags {
		if val, found = diagnosticExtra[T](diag.Extra); found {
			return
		}
	}
	return
}

type Diag hcl.Diagnostics // Diagnostics does implement error interface, but not, itself, an error.

func FindByExtra[T any](diags Diag) *hcl.Diagnostic {
	for _, diag := range diags {
		if _, found := diagnosticExtra[T](diag.Extra); found {
			return diag
		}
	}
	return nil
}

func (d Diag) Error() string {
	return hcl.Diagnostics(d).Error()
}

// Appends diag to diagnostics, returns true if the just-appended diagnostic is an error.
func (d *Diag) Append(diag *hcl.Diagnostic) (addedErrors bool) {
	if diag != nil {
		*d = append(*d, diag)
		return diag.Severity == hcl.DiagError
	}
	return false
}

// Add new diagnostic error.
func (d *Diag) Add(summary, detail string) {
	*d = append(*d, &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  summary,
		Detail:   detail,
	})
}

// Add new diagnostic warning.
func (d *Diag) AddWarn(summary, detail string) {
	*d = append(*d, &hcl.Diagnostic{
		Severity: hcl.DiagWarning,
		Summary:  summary,
		Detail:   detail,
	})
}

// Set the subject for diagnostics if it's not already specified.
func (d *Diag) DefaultSubject(rng *hcl.Range) {
	for _, diag := range *d {
		if diag.Subject == nil {
			diag.Subject = rng
		}
	}
}

// Adds an extra if no extra with the same type is present.
func AddExtraIfMissing[T any](d Diag, extra T) {
	for _, diag := range d {
		if _, found := diagnosticExtra[T](diag.Extra); !found {
			AddExtra(diag, extra)
		}
	}
}

// Adds extra without replacing existing extras.
func AddExtra(diag *hcl.Diagnostic, extra any) {
	switch extraT := diag.Extra.(type) {
	case nil:
		diag.Extra = extra
	case extraList:
		diag.Extra = extraList(append(extraT, extra))
	default:
		diag.Extra = extraList{extraT, extra}
	}
}

// Adds extra without replacing existing extras.
func (d Diag) AddExtra(extra any) {
	for _, diag := range d {
		AddExtra(diag, extra)
	}
}

// Appends all diags to diagnostics, returns true if the just-appended diagnostics contain an error.
func (d *Diag) Extend(diags []*hcl.Diagnostic) (haveAddedErrors bool) {
	*d = append(*d, diags...)
	return hcl.Diagnostics(diags).HasErrors()
}

// HasErrors returns true if the receiver contains any diagnostics of
// severity DiagError.
func (d Diag) HasErrors() bool {
	return hcl.Diagnostics(d).HasErrors()
}

// Creates diagnostic and appends it if err != nil.
func (d *Diag) AppendErr(err error, summary string) (haveAddedErrors bool) {
	// The body of the function is moved into `appendErr` to convince golang to inline the
	// `AppendErr`, making `err != nil` as cheap as usual.
	// Otherwise each AppendErr would waste a slow golang call just to check that err == nil and
	// return false
	haveAddedErrors = err != nil
	if haveAddedErrors {
		appendErr(d, err, summary)
	}
	return
}

// AppendErr and appendErr together can't be inlined. We're forbidding go from inlining
// appendErr into AppendErr and thus preventing AppendErr inlining.
//
//go:noinline
func appendErr(d *Diag, err error, summary string) {
	*d = append(*d, &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  summary,
		Detail:   err.Error(),
		Extra:    err,
	})
}

func FromHcl(diag *hcl.Diagnostic) Diag {
	if diag != nil {
		return Diag{diag}
	}
	return nil
}

func FromErr(err error, summary string) Diag {
	if err != nil {
		return Diag{{
			Severity: hcl.DiagError,
			Summary:  summary,
			Detail:   err.Error(),
			Extra:    err,
		}}
	}
	return nil
}

func FromErrSubj(err error, summary string, subject *hcl.Range) Diag {
	if err != nil {
		return Diag{{
			Severity: hcl.DiagError,
			Summary:  summary,
			Detail:   err.Error(),
			Extra:    err,
			Subject:  subject,
		}}
	}
	return nil
}
